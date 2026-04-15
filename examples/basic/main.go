// Package main demonstrates a basic open.mp gamemode written in Go.
package main

import (
	"context"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

// BasicGamemode is our example gamemode implementation.
// It embeds the omp base event handler to get default implementations for all events.
type BasicGamemode struct {
	omp.BaseEventHandler
	tickCount int
}

func (gm *BasicGamemode) OnLoad(ctx context.Context) error {
	_ = core.Log("[BasicGamemode] Gamemode loaded!")
	return nil
}

func (gm *BasicGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	_ = core.Log("[BasicGamemode] Player connected")
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "Welcome to the server!")
	return nil
}

func (gm *BasicGamemode) OnPlayerDisconnect(ctx context.Context, event *omp.PlayerDisconnectEvent) error {
	_ = core.Log("[BasicGamemode] Player disconnected")
	return nil
}

func (gm *BasicGamemode) OnPlayerSpawn(ctx context.Context, event *omp.PlayerSpawnEvent) error {
	_ = core.Log("[BasicGamemode] Player spawned")
	return nil
}

func (gm *BasicGamemode) OnPlayerText(ctx context.Context, event *omp.PlayerTextEvent) (bool, error) {
	_ = core.Log("[BasicGamemode] Chat: " + event.Text.Clone())
	return false, nil
}

func (gm *BasicGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	_ = core.Log("[BasicGamemode] Command: " + event.Command.Clone())
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	if event.Command.EqualString("/help") {
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "Commands: /help")
		return true, nil
	}
	return false, nil
}

func (gm *BasicGamemode) OnPlayerDeath(ctx context.Context, event *omp.PlayerDeathEvent) error {
	_ = core.Log("[BasicGamemode] Player died")
	return nil
}

func (gm *BasicGamemode) OnTick(ctx context.Context, event *omp.TickEvent) error {
	gm.tickCount++
	return nil
}

func NewGamemode() runtime.Gamemode {
	return &BasicGamemode{}
}
