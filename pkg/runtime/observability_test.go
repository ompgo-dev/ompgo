package runtime

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"
)

func TestErrorObserverSnapshot(t *testing.T) {
	t.Parallel()

	observer := NewErrorObserver(log.New(io.Discard, "", 0))

	observer.EventErrorHandler(context.Background(), "OnPlayerConnect", nil, errors.New("e1"))
	observer.EventErrorHandler(context.Background(), "OnPlayerConnect", nil, errors.New("e2"))
	observer.EventErrorHandler(context.Background(), "OnTick", nil, errors.New("e3"))
	observer.LifecycleErrorHandler(context.Background(), "Component.OnReady", errors.New("l1"))

	snap := observer.Snapshot()
	if snap.TotalEventErrors != 3 {
		t.Fatalf("TotalEventErrors=%d want=3", snap.TotalEventErrors)
	}
	if snap.TotalLifecycleErrors != 1 {
		t.Fatalf("TotalLifecycleErrors=%d want=1", snap.TotalLifecycleErrors)
	}
	if snap.EventCounts["OnPlayerConnect"] != 2 {
		t.Fatalf("OnPlayerConnect=%d want=2", snap.EventCounts["OnPlayerConnect"])
	}
	if snap.EventCounts["OnTick"] != 1 {
		t.Fatalf("OnTick=%d want=1", snap.EventCounts["OnTick"])
	}
	if snap.LifecycleCounts["Component.OnReady"] != 1 {
		t.Fatalf("Component.OnReady=%d want=1", snap.LifecycleCounts["Component.OnReady"])
	}
}
