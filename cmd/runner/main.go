package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/eva01/backup-guardian/domain"
	"github.com/eva01/backup-guardian/environment"
	"github.com/eva01/backup-guardian/migrations"
	"github.com/eva01/backup-guardian/runner"
	"github.com/eva01/backup-guardian/store"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	vars := environment.Parse()

	logger := newLogger(vars.LogLevel)

	if err := os.MkdirAll(filepath.Dir(vars.DBPath()), 0755); err != nil {
		log.Fatalf("could not create data directory: %v", err)
	}

	db, err := sql.Open("sqlite", vars.DBPath())
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	// Migrations (goose) — appliquées au démarrage
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("could not set dialect: %v", err)
	}
	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	s := store.New(store.WithDB(db))

	interval, err := vars.SyncIntervalDuration()
	if err != nil {
		log.Fatalf("invalid sync interval %q: %v", vars.SyncInterval, err)
	}

	r := runner.New(
		runner.WithStore(s.SyncRuns),
		runner.WithRcloneExecutor(&runner.LibraryRcloneExecutor{}),
		runner.WithScheduler(runner.NewScheduler(interval)),
		runner.WithSyncJob(&domain.SyncJob{
			Name:        "gdrive-to-s3",
			Source:      vars.SyncSource,
			Destination: vars.SyncDest,
		}),
		runner.WithLogger(logger),
	)

	if err := r.Run(context.Background(), vars); err != nil {
		log.Fatalf("runner failed: %v", err)
	}
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}
