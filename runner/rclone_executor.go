package runner

import (
	"context"
	"time"

	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/sync"

	"github.com/eva01/backup-guardian/runner/result"
)

// RcloneExecutor executes rclone sync operations.
type RcloneExecutor interface {
	Sync(ctx context.Context, source, dest string) (*result.RcloneResult, error)
}

// LibraryRcloneExecutor implements RcloneExecutor using the rclone Go library.
type LibraryRcloneExecutor struct{}

// Sync runs rclone sync from source to dest using the rclone library.
func (e *LibraryRcloneExecutor) Sync(ctx context.Context, source, dest string) (*result.RcloneResult, error) {
	start := time.Now()

	if err := fs.GlobalOptionsInit(); err != nil {
		return nil, err
	}

	fsrc, err := fs.NewFs(ctx, source)
	if err != nil {
		return &result.RcloneResult{Duration: time.Since(start)}, err
	}

	fdst, err := fs.NewFs(ctx, dest)
	if err != nil {
		return &result.RcloneResult{Duration: time.Since(start)}, err
	}

	if err := sync.Sync(ctx, fdst, fsrc, true); err != nil {
		return &result.RcloneResult{Duration: time.Since(start)}, err
	}

	return &result.RcloneResult{Duration: time.Since(start)}, nil
}
