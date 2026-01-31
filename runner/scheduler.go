package runner

import (
	"context"
	"time"
)

// Scheduler triggers syncs at regular intervals.
type Scheduler struct {
	interval time.Duration
	tick     chan struct{}
}

// NewScheduler creates a scheduler that fires at the given interval.
func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		interval: interval,
		tick:     make(chan struct{}, 1),
	}
}

// Run starts the scheduler loop, sending on C() at each interval.
// Stops when ctx is cancelled.
func (scheduler *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(scheduler.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case scheduler.tick <- struct{}{}:
			default:
			}
		}
	}
}

// C returns the channel that receives a signal at each scheduled tick.
func (scheduler *Scheduler) C() <-chan struct{} {
	return scheduler.tick
}
