package runtime

import (
	"context"
	"time"
)

type runtimeContextKey string

const (
	requestIDContextKey  runtimeContextKey = "ompgo.request_id"
	eventNameContextKey  runtimeContextKey = "ompgo.event_name"
	eventStartContextKey runtimeContextKey = "ompgo.event_started_at"
)

// WithRequestID stores a request identifier in context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

// RequestIDFromContext returns a request identifier if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDContextKey).(string)
	return v, ok && v != ""
}

// WithEventName stores the current event name in context.
func WithEventName(ctx context.Context, eventName string) context.Context {
	return context.WithValue(ctx, eventNameContextKey, eventName)
}

// EventNameFromContext returns the current event name if present.
func EventNameFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(eventNameContextKey).(string)
	return v, ok && v != ""
}

// WithEventStartedAt stores the event dispatch timestamp in context.
func WithEventStartedAt(ctx context.Context, startedAt time.Time) context.Context {
	return context.WithValue(ctx, eventStartContextKey, startedAt)
}

// EventStartedAtFromContext returns the event dispatch timestamp if present.
func EventStartedAtFromContext(ctx context.Context) (time.Time, bool) {
	v, ok := ctx.Value(eventStartContextKey).(time.Time)
	return v, ok && !v.IsZero()
}
