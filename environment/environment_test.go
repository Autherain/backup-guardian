package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariables_DBPath(t *testing.T) {
	v := &Variables{DataDir: "/data"}
	assert.Equal(t, "/data/backup-guardian.db", v.DBPath())

	v.DataDir = "."
	assert.Equal(t, "backup-guardian.db", v.DBPath())
}

func TestVariables_SyncIntervalDuration(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := &Variables{SyncInterval: "1h"}
		d, err := v.SyncIntervalDuration()
		require.NoError(t, err)
		assert.Equal(t, 3600, int(d.Seconds()))
	})

	t.Run("invalid", func(t *testing.T) {
		v := &Variables{SyncInterval: "not-a-duration"}
		_, err := v.SyncIntervalDuration()
		require.Error(t, err)
	})
}
