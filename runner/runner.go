package runner

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eva01/backup-guardian/domain"
	"github.com/eva01/backup-guardian/environment"
)

// Runner runs the backup sync loop.
type Runner struct {
	store     domain.SyncRunsReadWriter
	executor  RcloneExecutor
	scheduler *Scheduler
	job       *domain.SyncJob
	logger    *slog.Logger
}

// Option configures the runner.
type Option func(*Runner)

// New creates a new runner.
func New(options ...Option) *Runner {
	r := &Runner{
		logger: slog.Default(),
	}

	for _, opt := range options {
		opt(r)
	}

	return r
}

// WithStore sets the store.
func WithStore(store domain.SyncRunsReadWriter) Option {
	return func(r *Runner) { r.store = store }
}

// WithRcloneExecutor sets the rclone executor.
func WithRcloneExecutor(executor RcloneExecutor) Option {
	return func(r *Runner) { r.executor = executor }
}

// WithScheduler sets the scheduler.
func WithScheduler(scheduler *Scheduler) Option {
	return func(r *Runner) { r.scheduler = scheduler }
}

// WithSyncJob sets the sync job config.
func WithSyncJob(job *domain.SyncJob) Option {
	return func(r *Runner) { r.job = job }
}

// WithLogger sets the logger.
func WithLogger(logger *slog.Logger) Option {
	return func(r *Runner) { r.logger = logger }
}

// Run starts the runner loop. Blocks until context is cancelled or a signal is received.
func (r *Runner) Run(ctx context.Context, vars *environment.Variables) error {
	if r.store == nil {
		panic("runner requires store")
	}
	if r.executor == nil {
		panic("runner requires rclone executor")
	}
	if r.job == nil {
		panic("runner requires sync job")
	}

	interval, err := vars.SyncIntervalDuration()
	if err != nil {
		return err
	}

	if r.scheduler == nil {
		r.scheduler = NewScheduler(interval)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go r.scheduler.Run(runCtx)

	r.logger.Info("Runner started", slog.String("source", r.job.Source), slog.String("dest", r.job.Destination), slog.Duration("interval", interval))

	r.runSync(runCtx)

	for {
		select {
		case <-runCtx.Done():
			r.logger.Info("Runner stopping")
			return nil
		case sig := <-sigCh:
			r.logger.Info("Received signal, stopping", slog.String("signal", sig.String()))
			cancel()
			return nil
		case <-r.scheduler.C():
			r.runSync(runCtx)
		}
	}
}

func (r *Runner) runSync(ctx context.Context) {
	run := &domain.SyncRun{
		ID:        domain.NewSyncRunID(),
		JobName:   r.job.Name,
		Status:    domain.StatusRunning,
		StartedAt: time.Now(),
	}

	created, err := r.store.CreateSyncRun(run)
	if err != nil {
		r.logger.Error("Failed to create sync run", slog.Any("error", err))
		return
	}

	r.logger.Info("Starting sync", slog.String("run_id", created.ID), slog.String("job", r.job.Name))

	result, err := r.executor.Sync(ctx, r.job.Source, r.job.Destination)

	run = created
	run.FinishedAt = time.Now()
	run.FilesTransferred = 0
	run.BytesTransferred = 0

	if result != nil {
		run.FilesTransferred = result.FilesTransferred
		run.BytesTransferred = result.BytesTransferred
	}

	if err != nil {
		run.Status = domain.StatusFailed
		run.ErrorMessage = err.Error()
		r.logger.Error("Sync failed", slog.String("run_id", created.ID), slog.Any("error", err))
	} else {
		run.Status = domain.StatusSuccess
		r.logger.Info("Sync completed", slog.String("run_id", created.ID),
			slog.Int64("files", run.FilesTransferred),
			slog.Int64("bytes", run.BytesTransferred),
			slog.Duration("duration", result.Duration))
	}

	if updateErr := r.store.UpdateSyncRun(run); updateErr != nil {
		r.logger.Error("Failed to update sync run", slog.String("run_id", created.ID), slog.Any("error", updateErr))
	}
}
