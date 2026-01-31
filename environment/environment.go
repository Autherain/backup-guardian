// Package environment provides methods to interact with environment variables.
// Please, make sure to update the .env.example file when modifying this structure.
package environment

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Variables represents the environment variables used by the application.
type Variables struct {
	// DataDir is the base directory for persistent data (DB, etc.). Mount this in Docker.
	DataDir string `env:"BG_DATA_DIR" envDefault:"."`

	// Remote names must match sections in rclone.conf (see rclone.conf.example).
	SyncSource string `env:"BG_SYNC_SOURCE,required" envDefault:"gdrive:"`
	SyncDest   string `env:"BG_SYNC_DEST,required" envDefault:"s3:bucket-name/backups"`
	SyncInterval string `env:"BG_SYNC_INTERVAL" envDefault:"6h"`

	LogLevel string `env:"BG_LOG_LEVEL" envDefault:"info"`
}

// DBPath returns the SQLite database path, derived from DataDir.
func (v *Variables) DBPath() string {
	return filepath.Join(v.DataDir, "backup-guardian.db")
}

// SyncIntervalDuration returns the parsed sync interval.
func (v *Variables) SyncIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(v.SyncInterval)
}

// Parse environment variables.
func Parse() *Variables {
	godotenv.Load() // Used for local development.

	result := &Variables{}
	if err := env.Parse(result); err != nil {
		panic(fmt.Errorf("could not parse environment variables: %w", err))
	}

	return result
}
