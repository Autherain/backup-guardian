package domain

import (
	"time"

	"github.com/eva01/backup-guardian/internal/errors"
	"github.com/google/uuid"
)

const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

// SyncRun represents a single sync execution.
type SyncRun struct {
	ID               string
	JobName          string
	Status           string
	StartedAt        time.Time
	FinishedAt       time.Time
	ErrorMessage     string
	FilesTransferred int64
	BytesTransferred int64
	CreatedAt        time.Time
}

// SyncRunSelector identifies a sync run for reads.
type SyncRunSelector struct {
	ID string
}

// SyncRunsSelector filters sync runs for listing.
type SyncRunsSelector struct {
	JobName string
	Limit   int
	Offset  int
}

// Validate validates the sync run.
func (r *SyncRun) Validate() error {
	if r.ID == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "ID must be set"}
	}
	if r.JobName == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "JobName must be set"}
	}
	if r.Status == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Status must be set"}
	}

	return nil
}

// NewSyncRunID returns a new UUID for a sync run.
func NewSyncRunID() string {
	return uuid.New().String()
}

// SyncRunsReadWriter combines read and write operations for sync runs.
type SyncRunsReadWriter interface {
	SyncRunsReader
	SyncRunsWriter
}

// SyncRunsReader defines read operations.
type SyncRunsReader interface {
	GetSyncRun(selector *SyncRunSelector) (*SyncRun, error)
	ListSyncRuns(selector *SyncRunsSelector) ([]*SyncRun, error)
}

// SyncRunsWriter defines write operations.
type SyncRunsWriter interface {
	CreateSyncRun(run *SyncRun) (*SyncRun, error)
	UpdateSyncRun(run *SyncRun) error
}
