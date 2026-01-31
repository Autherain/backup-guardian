package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncRun_Validate(t *testing.T) {
	validRun := &SyncRun{
		ID:       "run-1",
		JobName:  "job-1",
		Status:   StatusRunning,
	}

	t.Run("valid", func(t *testing.T) {
		err := validRun.Validate()
		require.NoError(t, err)
	})

	t.Run("empty ID", func(t *testing.T) {
		r := &SyncRun{ID: "", JobName: "job", Status: StatusPending}
		err := r.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID must be set")
	})

	t.Run("empty JobName", func(t *testing.T) {
		r := &SyncRun{ID: "id", JobName: "", Status: StatusPending}
		err := r.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "JobName must be set")
	})

	t.Run("empty Status", func(t *testing.T) {
		r := &SyncRun{ID: "id", JobName: "job", Status: ""}
		err := r.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Status must be set")
	})
}

func TestNewSyncRunID(t *testing.T) {
	id := NewSyncRunID()
	require.NotEmpty(t, id)
	_, err := uuid.Parse(id)
	require.NoError(t, err)
}
