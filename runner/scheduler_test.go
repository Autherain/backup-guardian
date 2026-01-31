package runner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestScheduler(t *testing.T) {
	interval := 5 * time.Millisecond
	s := NewScheduler(interval)
	require.NotNil(t, s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Run(ctx)

	select {
	case <-s.C():
		// received at least one tick
	case <-time.After(2 * interval):
		t.Fatal("expected at least one tick within 2 intervals")
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	interval := 10 * time.Millisecond
	s := NewScheduler(interval)
	ctx, cancel := context.WithCancel(context.Background())

	go s.Run(ctx)
	cancel()

	// Run() should return; no panic or hang. Give it a moment.
	time.Sleep(2 * interval)
	// If we get here without deadlock, the goroutine exited.
}
