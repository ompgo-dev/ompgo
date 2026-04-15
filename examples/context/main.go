// Package main demonstrates context propagation and error handling patterns.
package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

var eventSeq atomic.Uint64

type ContextGamemode struct {
	omp.BaseEventHandler
}

func (gm *ContextGamemode) OnLoad(ctx context.Context) error {
	_ = core.Log("[ContextDemo] Loaded")
	return nil
}

func (gm *ContextGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}

	name, _ := players.GetName(event.Player)
	reqID := readRequestID(ctx)
	_ = players.SendClientMessage(
		event.Player,
		uint32(omp.ColorGreen),
		fmt.Sprintf("Welcome %s! request_id=%s", name, reqID),
	)
	_ = players.SendClientMessage(
		event.Player,
		uint32(omp.ColorYellow),
		"Try: /ctx, /cancel, /fail",
	)
	return nil
}

func (gm *ContextGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}

	command := event.Command.Clone()

	switch command {
	case "/ctx":
		reqID := readRequestID(ctx)
		age := time.Since(readStartedAt(ctx)).Truncate(time.Millisecond)
		eventName, _ := runtime.EventNameFromContext(ctx)
		_ = players.SendClientMessage(
			event.Player,
			uint32(omp.ColorWhite),
			fmt.Sprintf("ctx event=%s request_id=%s age=%s", eventName, reqID, age),
		)
		return true, nil

	case "/cancel":
		child, cancel := context.WithCancel(ctx)
		cancel()
		select {
		case <-child.Done():
			_ = players.SendClientMessage(event.Player, uint32(omp.ColorOrange), "child context cancelled")
		default:
			_ = players.SendClientMessage(event.Player, uint32(omp.ColorRed), "child context still active")
		}
		return true, nil

	case "/fail":
		// Intentional error for demonstrating WithEventErrorHandler.
		return true, fmt.Errorf("intentional failure from /fail")
	}

	return false, nil
}

func NewGamemode() runtime.Gamemode {
	return &ContextGamemode{}
}

func init() {
	runtime.Bootstrap(
		runtime.WithComponentName("ompgo_context_demo"),
		runtime.WithGamemode(NewGamemode),
		runtime.WithSetup(func(ctx context.Context) error {
			log.Printf("[ContextDemo] setup complete")
			return nil
		}),
		runtime.WithOnReady(func(ctx context.Context) error {
			log.Printf("[ContextDemo] component ready")
			return nil
		}),
		runtime.WithOnFree(func(ctx context.Context) error {
			log.Printf("[ContextDemo] component free")
			return nil
		}),
		runtime.WithEventErrorPolicy(runtime.ErrorPolicyBlockOnError),
		runtime.WithLifecycleErrorPolicy(runtime.ErrorPolicyContinue),
		runtime.WithEventContextDecorator(func(ctx context.Context, eventName string, event any) context.Context {
			_ = event
			seq := eventSeq.Add(1)
			return runtime.WithRequestID(ctx, fmt.Sprintf("%s-%d", eventName, seq))
		}),
		runtime.WithEventErrorHandler(func(ctx context.Context, eventName string, event any, err error) {
			_ = ctx
			_ = event
			log.Printf("[ContextDemo] event error in %s: %v", eventName, err)
		}),
		runtime.WithLifecycleErrorHandler(func(ctx context.Context, stage string, err error) {
			_ = ctx
			log.Printf("[ContextDemo] lifecycle error in %s: %v", stage, err)
		}),
	)
}

func readRequestID(ctx context.Context) string {
	if v, ok := runtime.RequestIDFromContext(ctx); ok && v != "" {
		return v
	}
	return "unknown"
}

func readStartedAt(ctx context.Context) time.Time {
	if v, ok := runtime.EventStartedAtFromContext(ctx); ok {
		return v
	}
	return time.Now()
}

func main() {}
