//go:build integration

package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLibraryRcloneExecutor_Sync_Integration_Local(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a file in source so there is something to sync
	err := os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	e := &LibraryRcloneExecutor{}
	ctx := context.Background()
	// Local paths: pass absolute path; rclone treats it as local without config
	source := srcDir
	dest := dstDir

	syncResult, err := e.Sync(ctx, source, dest)
	require.NoError(t, err)
	require.NotNil(t, syncResult)
	require.GreaterOrEqual(t, syncResult.Duration, time.Duration(0))

	// Verify file was synced to destination
	destPath := filepath.Join(dstDir, "test.txt")
	_, err = os.Stat(destPath)
	require.NoError(t, err)
}

// TestLibraryRcloneExecutor_Sync_Integration_Memory uses rclone's :memory: backend
// (in-RAM, no config, no disk). Syncs local -> :memory:src -> :memory:dst -> local
// and verifies the file is present.
func TestLibraryRcloneExecutor_Sync_Integration_Memory(t *testing.T) {
	localSrc := t.TempDir()
	localDst := t.TempDir()

	err := os.WriteFile(filepath.Join(localSrc, "test.txt"), []byte("hello memory"), 0644)
	require.NoError(t, err)

	e := &LibraryRcloneExecutor{}
	ctx := context.Background()

	// Populate :memory:src from local (path only = local backend)
	_, err = e.Sync(ctx, localSrc, ":memory:src")
	require.NoError(t, err)

	// Sync between two memory remotes (no disk, no credentials)
	_, err = e.Sync(ctx, ":memory:src", ":memory:dst")
	require.NoError(t, err)

	// Pull back to local to verify content
	_, err = e.Sync(ctx, ":memory:dst", localDst)
	require.NoError(t, err)

	destPath := filepath.Join(localDst, "test.txt")
	content, err := os.ReadFile(destPath)
	require.NoError(t, err)
	require.Equal(t, "hello memory", string(content))
}
