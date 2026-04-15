package runtime

import (
	"context"
	"testing"
	"time"
)

func TestNewEventContextPreservesBaseValuesAndMetadata(t *testing.T) {
	oldCfg := *currentEventDispatchConfig()
	t.Cleanup(func() { setEventDispatchConfig(oldCfg) })

	cfg := oldCfg
	cfg.baseContext = WithRequestID(context.Background(), "req-123")
	cfg.contextProvider = nil
	cfg.decorator = nil
	setEventDispatchConfig(cfg)

	ctx := newEventContext(currentEventDispatchConfig(), "OnTick")

	if reqID, ok := RequestIDFromContext(ctx); !ok || reqID != "req-123" {
		t.Fatalf("RequestIDFromContext = (%q, %v), want (req-123, true)", reqID, ok)
	}
	if name, ok := EventNameFromContext(ctx); !ok || name != "OnTick" {
		t.Fatalf("EventNameFromContext = (%q, %v), want (OnTick, true)", name, ok)
	}
	startedAt, ok := EventStartedAtFromContext(ctx)
	if !ok {
		t.Fatal("EventStartedAtFromContext() = (_, false), want true")
	}
	if time.Since(startedAt) > time.Second {
		t.Fatalf("EventStartedAtFromContext() = %v, want recent timestamp", startedAt)
	}
}

func TestNewEventContextDecoratorSeesEventMetadata(t *testing.T) {
	oldCfg := *currentEventDispatchConfig()
	t.Cleanup(func() { setEventDispatchConfig(oldCfg) })

	decorated := false
	cfg := oldCfg
	cfg.decorator = func(ctx context.Context, eventName string, event any) context.Context {
		if name, ok := EventNameFromContext(ctx); !ok || name != eventName {
			t.Fatalf("decorator EventNameFromContext = (%q, %v), want (%q, true)", name, ok, eventName)
		}
		if _, ok := EventStartedAtFromContext(ctx); !ok {
			t.Fatal("decorator missing started-at metadata")
		}
		decorated = true
		return WithRequestID(ctx, "decorated")
	}
	setEventDispatchConfig(cfg)

	cfgPtr := currentEventDispatchConfig()
	ctx := newEventContext(cfgPtr, "OnPlayerConnect")
	ctx = decorateEventContext(cfgPtr, ctx, "OnPlayerConnect", struct{}{})

	if !decorated {
		t.Fatal("decorator was not invoked")
	}
	if reqID, ok := RequestIDFromContext(ctx); !ok || reqID != "decorated" {
		t.Fatalf("RequestIDFromContext = (%q, %v), want (decorated, true)", reqID, ok)
	}
}
