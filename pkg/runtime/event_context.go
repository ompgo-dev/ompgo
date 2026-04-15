package runtime

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// EventContextDecorator can decorate the context used for a specific event dispatch.
type EventContextDecorator func(ctx context.Context, eventName string, event any) context.Context

// EventErrorHandler receives errors returned by event handlers and gamemode handlers.
type EventErrorHandler func(ctx context.Context, eventName string, event any, err error)

// LifecycleErrorHandler receives errors returned by lifecycle hooks.
type LifecycleErrorHandler func(ctx context.Context, stage string, err error)

// ErrorPolicy controls runtime behavior after callback errors.
type ErrorPolicy int

const (
	// ErrorPolicyContinue reports errors and continues execution.
	ErrorPolicyContinue ErrorPolicy = iota
	// ErrorPolicyBlockOnError reports errors and blocks further processing.
	ErrorPolicyBlockOnError
)

type eventDispatchConfig struct {
	baseContext           context.Context
	contextProvider       func() context.Context
	decorator             EventContextDecorator
	errorHandler          EventErrorHandler
	lifecycleErrorHandler LifecycleErrorHandler
	eventErrorPolicy      ErrorPolicy
	lifecycleErrorPolicy  ErrorPolicy
}

var (
	defaultEventDispatch = &eventDispatchConfig{
		baseContext:          context.Background(),
		eventErrorPolicy:     ErrorPolicyContinue,
		lifecycleErrorPolicy: ErrorPolicyContinue,
	}
	eventDispatch atomic.Pointer[eventDispatchConfig]
)

type eventContext struct {
	base      context.Context
	eventName string
	startedAt time.Time
}

func (ctx eventContext) Deadline() (time.Time, bool) {
	return ctx.base.Deadline()
}

func (ctx eventContext) Done() <-chan struct{} {
	return ctx.base.Done()
}

func (ctx eventContext) Err() error {
	return ctx.base.Err()
}

func (ctx eventContext) Value(key any) any {
	switch key {
	case eventNameContextKey:
		return ctx.eventName
	case eventStartContextKey:
		return ctx.startedAt
	default:
		return ctx.base.Value(key)
	}
}

func setEventDispatchConfig(cfg eventDispatchConfig) {
	if cfg.baseContext == nil {
		cfg.baseContext = context.Background()
	}
	if cfg.eventErrorPolicy < ErrorPolicyContinue || cfg.eventErrorPolicy > ErrorPolicyBlockOnError {
		cfg.eventErrorPolicy = ErrorPolicyContinue
	}
	if cfg.lifecycleErrorPolicy < ErrorPolicyContinue || cfg.lifecycleErrorPolicy > ErrorPolicyBlockOnError {
		cfg.lifecycleErrorPolicy = ErrorPolicyContinue
	}
	next := cfg
	eventDispatch.Store(&next)
}

func currentEventDispatchConfig() *eventDispatchConfig {
	if cfg := eventDispatch.Load(); cfg != nil {
		return cfg
	}
	return defaultEventDispatch
}

func newEventContext(cfg *eventDispatchConfig, eventName string) context.Context {
	base := resolveBaseContext(cfg)

	return context.Context(eventContext{
		base:      base,
		eventName: eventName,
		startedAt: time.Now(),
	})
}

func decorateEventContext(cfg *eventDispatchConfig, ctx context.Context, eventName string, event any) context.Context {
	if cfg.decorator != nil {
		if decorated := cfg.decorator(ctx, eventName, event); decorated != nil {
			return decorated
		}
	}
	return ctx
}

func newLifecycleContext(stage string) (context.Context, context.CancelFunc) {
	cfg := currentEventDispatchConfig()
	base := resolveBaseContext(cfg)
	return context.WithCancel(base)
}

func reportEventErrorWithConfig(ctx context.Context, cfg *eventDispatchConfig, eventName string, event any, err error) {
	if err == nil {
		return
	}
	if cfg.errorHandler != nil {
		cfg.errorHandler(ctx, eventName, event, err)
		return
	}
	fmt.Printf("[Runtime] Event handler error in %s: %v\n", eventName, err)
}

func reportEventError(ctx context.Context, eventName string, event any, err error) {
	reportEventErrorWithConfig(ctx, currentEventDispatchConfig(), eventName, event, err)
}

func reportLifecycleError(ctx context.Context, stage string, err error) {
	if err == nil {
		return
	}
	cfg := currentEventDispatchConfig()
	if cfg.lifecycleErrorHandler != nil {
		cfg.lifecycleErrorHandler(ctx, stage, err)
		return
	}
	fmt.Printf("[Runtime] Lifecycle error in %s: %v\n", stage, err)
}

func shouldBlockOnEventError() bool {
	return eventErrorPolicyBlocks(currentEventDispatchConfig())
}

func eventErrorPolicyBlocks(cfg *eventDispatchConfig) bool {
	return cfg.eventErrorPolicy == ErrorPolicyBlockOnError
}

func shouldBlockOnLifecycleError() bool {
	cfg := currentEventDispatchConfig()
	return cfg.lifecycleErrorPolicy == ErrorPolicyBlockOnError
}

func resolveBaseContext(cfg *eventDispatchConfig) context.Context {
	base := cfg.baseContext
	if cfg.contextProvider != nil {
		if provided := cfg.contextProvider(); provided != nil {
			base = provided
		}
	}
	if base == nil {
		base = context.Background()
	}
	return base
}
