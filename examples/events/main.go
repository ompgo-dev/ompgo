// Package main demonstrates event handling using the current API.
// Build with: CGO_ENABLED=1 go build -buildmode=plugin -o events.so .
package main

import (
	"context"
	"fmt"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

type EventsGamemode struct {
	omp.BaseEventHandler
}

func (g *EventsGamemode) OnLoad(ctx context.Context) error {
	_ = core.Log("[Events] Gamemode loaded")
	return nil
}

func (g *EventsGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	name, _ := players.GetName(event.Player)
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), fmt.Sprintf("Welcome, %s!", name))
	return nil
}

func (g *EventsGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	if event.Command.EqualString("/ping") {
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "pong")
		return true, nil
	}
	return false, nil
}

func NewGamemode() runtime.Gamemode {
	return &EventsGamemode{}
}
