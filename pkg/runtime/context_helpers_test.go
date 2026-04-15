package runtime

import (
	"context"
	"testing"
	"time"
)

func TestContextHelpers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithEventName(ctx, "OnPlayerConnect")
	now := time.Now().UTC()
	ctx = WithEventStartedAt(ctx, now)

	if reqID, ok := RequestIDFromContext(ctx); !ok || reqID != "req-123" {
		t.Fatalf("RequestIDFromContext = (%q, %v), want (req-123, true)", reqID, ok)
	}
	if name, ok := EventNameFromContext(ctx); !ok || name != "OnPlayerConnect" {
		t.Fatalf("EventNameFromContext = (%q, %v), want (OnPlayerConnect, true)", name, ok)
	}
	if startedAt, ok := EventStartedAtFromContext(ctx); !ok || !startedAt.Equal(now) {
		t.Fatalf("EventStartedAtFromContext = (%v, %v), want (%v, true)", startedAt, ok, now)
	}
}
