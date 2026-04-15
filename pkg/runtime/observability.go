package runtime

import (
	"context"
	"log"
	"sync"
)

// ErrorSnapshot contains aggregated error counts.
type ErrorSnapshot struct {
	TotalEventErrors     uint64
	TotalLifecycleErrors uint64
	EventCounts          map[string]uint64
	LifecycleCounts      map[string]uint64
}

// ErrorObserver aggregates runtime errors and optionally logs them.
type ErrorObserver struct {
	logger *log.Logger

	mu                   sync.Mutex
	totalEventErrors     uint64
	totalLifecycleErrors uint64
	eventCounts          map[string]uint64
	lifecycleCounts      map[string]uint64
}

// NewErrorObserver creates a new observer. Pass nil logger to disable logging.
func NewErrorObserver(logger *log.Logger) *ErrorObserver {
	return &ErrorObserver{
		logger:          logger,
		eventCounts:     map[string]uint64{},
		lifecycleCounts: map[string]uint64{},
	}
}

// EventErrorHandler is suitable for runtime.WithEventErrorHandler.
func (o *ErrorObserver) EventErrorHandler(ctx context.Context, eventName string, event any, err error) {
	_ = ctx
	_ = event
	if err == nil {
		return
	}
	o.mu.Lock()
	o.totalEventErrors++
	o.eventCounts[eventName]++
	o.mu.Unlock()
	if o.logger != nil {
		o.logger.Printf("[runtime] event error event=%s err=%v", eventName, err)
	}
}

// LifecycleErrorHandler is suitable for runtime.WithLifecycleErrorHandler.
func (o *ErrorObserver) LifecycleErrorHandler(ctx context.Context, stage string, err error) {
	_ = ctx
	if err == nil {
		return
	}
	o.mu.Lock()
	o.totalLifecycleErrors++
	o.lifecycleCounts[stage]++
	o.mu.Unlock()
	if o.logger != nil {
		o.logger.Printf("[runtime] lifecycle error stage=%s err=%v", stage, err)
	}
}

// Snapshot returns a copy of current counts.
func (o *ErrorObserver) Snapshot() ErrorSnapshot {
	o.mu.Lock()
	defer o.mu.Unlock()

	eventCounts := make(map[string]uint64, len(o.eventCounts))
	for k, v := range o.eventCounts {
		eventCounts[k] = v
	}
	lifecycleCounts := make(map[string]uint64, len(o.lifecycleCounts))
	for k, v := range o.lifecycleCounts {
		lifecycleCounts[k] = v
	}

	return ErrorSnapshot{
		TotalEventErrors:     o.totalEventErrors,
		TotalLifecycleErrors: o.totalLifecycleErrors,
		EventCounts:          eventCounts,
		LifecycleCounts:      lifecycleCounts,
	}
}

// WithErrorObserver wires an observer into both event and lifecycle error pipelines.
func WithErrorObserver(observer *ErrorObserver) Option {
	return func(cfg *Config) {
		if observer == nil {
			return
		}

		prevEvent := cfg.EventError
		cfg.EventError = func(ctx context.Context, eventName string, event any, err error) {
			observer.EventErrorHandler(ctx, eventName, event, err)
			if prevEvent != nil {
				prevEvent(ctx, eventName, event, err)
			}
		}

		prevLifecycle := cfg.LifecycleError
		cfg.LifecycleError = func(ctx context.Context, stage string, err error) {
			observer.LifecycleErrorHandler(ctx, stage, err)
			if prevLifecycle != nil {
				prevLifecycle(ctx, stage, err)
			}
		}
	}
}
