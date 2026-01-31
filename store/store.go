package store

import (
	"database/sql"
	"fmt"

	"github.com/eva01/backup-guardian/domain"
)

// Store provides access to persistence layers.
type Store struct {
	SyncRuns domain.SyncRunsReadWriter

	db *sql.DB
}

// Option configures the store.
type Option func(*Store) error

// New returns a store configured with the given options.
func New(options ...Option) *Store {
	s := &Store{}

	s.SyncRuns = &syncRunsStore{baseStore: s}

	for _, opt := range options {
		if err := opt(s); err != nil {
			panic(fmt.Errorf("could not create store: %w", err))
		}
	}

	return s
}

// WithDB sets the database connection.
func WithDB(db *sql.DB) Option {
	return func(s *Store) error {
		if err := db.Ping(); err != nil {
			return fmt.Errorf("could not ping database: %w", err)
		}
		s.db = db
		return nil
	}
}
