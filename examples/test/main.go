// Package main implements a simple test gamemode to verify the ompgo integration.
package main

import (
	"context"
	"log"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/all"
	"github.com/ompgo-dev/ompgo/pkg/omp/class"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

// TestGamemode is a minimal gamemode for testing
type TestGamemode struct {
	omp.BaseEventHandler
}

func (g *TestGamemode) OnLoad(ctx context.Context) error {
	log.Println("[TestGamemode] OnLoad called")

	_ = core.GameModeSetText("ompgo Test Mode")
	classID1 := int32(-1)
	classID2 := int32(-1)
	_ = class.Add(0, 0, 1958.3783, 1343.1572, 15.3746, 270.1425, 0, 0, 0, 0, 0, 0, &classID1)
	_ = class.Add(1, 0, 1958.3783, 1343.1572, 15.3746, 270.1425, 0, 0, 0, 0, 0, 0, &classID2)

	log.Println("[TestGamemode] Initialization complete")
	log.Println("[TestGamemode] Gamemode: ompgo Test Mode")
	log.Println("[TestGamemode] Version: 0.1.0")
	log.Println("[TestGamemode] Ready to accept players!")
	return nil
}

// OnPlayerConnect is called when a player connects
func (g *TestGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	log.Printf("[TestGamemode] Player connected")
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	name, _ := players.GetName(event.Player)
	log.Printf("[TestGamemode] Player name: %s", name)
	msg := "Welcome to ompgo Test Server, " + name + "!"
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), msg)
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "This server is running Go!")
	_ = all.SendClientMessage(uint32(omp.ColorGrey), name+" has joined the server.")
	return nil
}

// OnPlayerDisconnect is called when a player disconnects
func (g *TestGamemode) OnPlayerDisconnect(ctx context.Context, event *omp.PlayerDisconnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	name, _ := players.GetName(event.Player)
	log.Printf("[TestGamemode] Player (%s) disconnected", name)
	_ = all.SendClientMessage(uint32(omp.ColorGrey), name+" has left the server.")
	return nil
}

// OnPlayerSpawn is called when a player spawns
func (g *TestGamemode) OnPlayerSpawn(ctx context.Context, event *omp.PlayerSpawnEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	name, _ := players.GetName(event.Player)
	log.Printf("[TestGamemode] Player (%s) spawned", name)
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "You have spawned! Type /help for commands.")
	return nil
}

// OnPlayerText is called when a player sends a chat message
func (g *TestGamemode) OnPlayerText(ctx context.Context, event *omp.PlayerTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	name, _ := players.GetName(event.Player)
	log.Printf("[TestGamemode] %s: %s", name, event.Text.Clone())
	return false, nil
}

// OnPlayerCommandText is called when a player types a command
func (g *TestGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	name, _ := players.GetName(event.Player)
	log.Printf("[TestGamemode] %s typed: %s", name, event.Command.Clone())
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "Commands are logged to console")

	return true, nil
}

// NewGamemode returns the test gamemode instance
func NewGamemode() runtime.Gamemode {
	return &TestGamemode{}
}
