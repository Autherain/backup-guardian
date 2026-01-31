package runner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eva01/backup-guardian/domain"
	domainmocks "github.com/eva01/backup-guardian/domain/mocks"
	"github.com/eva01/backup-guardian/environment"
	"github.com/eva01/backup-guardian/runner"
	runnermocks "github.com/eva01/backup-guardian/runner/mocks"
	"github.com/eva01/backup-guardian/runner/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunner_Run_Success(t *testing.T) {
	storeMock := domainmocks.NewSyncRunsReadWriter(t)
	execMock := runnermocks.NewRcloneExecutor(t)

	createdRun := &domain.SyncRun{
		ID:        "test-run-id",
		JobName:   "test-job",
		Status:    domain.StatusRunning,
		StartedAt: time.Now(),
	}
	storeMock.On("CreateSyncRun", mock.Anything).Run(func(args mock.Arguments) {
		run := args.Get(0).(*domain.SyncRun)
		require.NotEmpty(t, run.ID)
		assert.Equal(t, "test-job", run.JobName)
		assert.Equal(t, domain.StatusRunning, run.Status)
	}).Return(createdRun, nil).Once()

	syncResult := &result.RcloneResult{
		Duration:         time.Second,
		FilesTransferred: 10,
		BytesTransferred: 100,
	}
	execMock.On("Sync", mock.Anything, "source", "dest").Return(syncResult, nil).Once()

	syncDone := make(chan struct{})
	storeMock.On("UpdateSyncRun", mock.Anything).Run(func(args mock.Arguments) {
		run := args.Get(0).(*domain.SyncRun)
		assert.Equal(t, "test-run-id", run.ID)
		assert.Equal(t, domain.StatusSuccess, run.Status)
		assert.Equal(t, int64(10), run.FilesTransferred)
		assert.Equal(t, int64(100), run.BytesTransferred)
		assert.Empty(t, run.ErrorMessage)
		close(syncDone)
	}).Return(nil).Once()

	vars := &environment.Variables{SyncInterval: "24h"}
	job := &domain.SyncJob{Name: "test-job", Source: "source", Destination: "dest"}

	r := runner.New(
		runner.WithStore(storeMock),
		runner.WithRcloneExecutor(execMock),
		runner.WithSyncJob(job),
		runner.WithScheduler(runner.NewScheduler(24*time.Hour)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- r.Run(ctx, vars) }()

	<-syncDone
	cancel()
	err := <-errCh
	require.NoError(t, err)
}

func TestRunner_Run_SyncFails(t *testing.T) {
	storeMock := domainmocks.NewSyncRunsReadWriter(t)
	execMock := runnermocks.NewRcloneExecutor(t)

	createdRun := &domain.SyncRun{
		ID:        "test-run-id",
		JobName:   "test-job",
		Status:    domain.StatusRunning,
		StartedAt: time.Now(),
	}
	storeMock.On("CreateSyncRun", mock.Anything).Return(createdRun, nil).Once()

	syncErr := errors.New("sync failed")
	execMock.On("Sync", mock.Anything, "source", "dest").Return(nil, syncErr).Once()

	syncDone := make(chan struct{})
	storeMock.On("UpdateSyncRun", mock.Anything).Run(func(args mock.Arguments) {
		run := args.Get(0).(*domain.SyncRun)
		assert.Equal(t, domain.StatusFailed, run.Status)
		assert.Equal(t, "sync failed", run.ErrorMessage)
		close(syncDone)
	}).Return(nil).Once()

	vars := &environment.Variables{SyncInterval: "24h"}
	job := &domain.SyncJob{Name: "test-job", Source: "source", Destination: "dest"}

	r := runner.New(
		runner.WithStore(storeMock),
		runner.WithRcloneExecutor(execMock),
		runner.WithSyncJob(job),
		runner.WithScheduler(runner.NewScheduler(24*time.Hour)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- r.Run(ctx, vars) }()

	<-syncDone
	cancel()
	err := <-errCh
	require.NoError(t, err)
}

func TestRunner_Run_CreateSyncRunFails(t *testing.T) {
	storeMock := domainmocks.NewSyncRunsReadWriter(t)
	execMock := runnermocks.NewRcloneExecutor(t)

	createErr := errors.New("db unavailable")
	createCalled := make(chan struct{})
	storeMock.On("CreateSyncRun", mock.Anything).Run(func(args mock.Arguments) {
		close(createCalled)
	}).Return(nil, createErr).Once()
	// Sync et UpdateSyncRun ne doivent jamais être appelés

	vars := &environment.Variables{SyncInterval: "24h"}
	job := &domain.SyncJob{Name: "test-job", Source: "source", Destination: "dest"}

	r := runner.New(
		runner.WithStore(storeMock),
		runner.WithRcloneExecutor(execMock),
		runner.WithSyncJob(job),
		runner.WithScheduler(runner.NewScheduler(24*time.Hour)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- r.Run(ctx, vars) }()

	<-createCalled
	cancel()
	err := <-errCh
	require.NoError(t, err)
}

func TestRunner_Run_UpdateSyncRunFails(t *testing.T) {
	storeMock := domainmocks.NewSyncRunsReadWriter(t)
	execMock := runnermocks.NewRcloneExecutor(t)

	createdRun := &domain.SyncRun{
		ID:        "test-run-id",
		JobName:   "test-job",
		Status:    domain.StatusRunning,
		StartedAt: time.Now(),
	}
	storeMock.On("CreateSyncRun", mock.Anything).Return(createdRun, nil).Once()

	syncResult := &result.RcloneResult{FilesTransferred: 5, BytesTransferred: 50}
	execMock.On("Sync", mock.Anything, "source", "dest").Return(syncResult, nil).Once()

	updateErr := errors.New("db write failed")
	syncDone := make(chan struct{})
	storeMock.On("UpdateSyncRun", mock.Anything).Run(func(args mock.Arguments) {
		run := args.Get(0).(*domain.SyncRun)
		assert.Equal(t, domain.StatusSuccess, run.Status)
		close(syncDone)
	}).Return(updateErr).Once()

	vars := &environment.Variables{SyncInterval: "24h"}
	job := &domain.SyncJob{Name: "test-job", Source: "source", Destination: "dest"}

	r := runner.New(
		runner.WithStore(storeMock),
		runner.WithRcloneExecutor(execMock),
		runner.WithSyncJob(job),
		runner.WithScheduler(runner.NewScheduler(24*time.Hour)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- r.Run(ctx, vars) }()

	<-syncDone
	cancel()
	err := <-errCh
	require.NoError(t, err)
}

func TestRunner_Run_InvalidSyncInterval(t *testing.T) {
	storeMock := domainmocks.NewSyncRunsReadWriter(t)
	execMock := runnermocks.NewRcloneExecutor(t)

	vars := &environment.Variables{SyncInterval: "invalid"}
	job := &domain.SyncJob{Name: "test-job", Source: "source", Destination: "dest"}

	r := runner.New(
		runner.WithStore(storeMock),
		runner.WithRcloneExecutor(execMock),
		runner.WithSyncJob(job),
	)

	ctx := context.Background()
	err := r.Run(ctx, vars)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}
