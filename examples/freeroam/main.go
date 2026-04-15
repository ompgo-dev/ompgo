// Package main provides a freeroam component written in Go.
package main

import (
	"context"
	"fmt"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

import "C"

type FreeroamGamemode struct {
	omp.BaseEventHandler
	tickCount int
	startTime int64
}

// OnLoad initializes the freeroam gamemode.
func (gm *FreeroamGamemode) OnLoad(ctx context.Context) error {
	gm.startTime = 0
	_ = core.Log("[Freeroam] Gamemode loaded!")
	_ = core.Log("[Freeroam] Players can use /help for available commands")
	return nil
}

func (gm *FreeroamGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	_ = core.Log("[Freeroam] Player connected")
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "Welcome to Freeroam!")
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "Type /help for commands")
	return nil
}

func (gm *FreeroamGamemode) OnPlayerSpawn(ctx context.Context, event *omp.PlayerSpawnEvent) error {
	_ = core.Log("[Freeroam] Player spawned")
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	_ = players.SetPos(event.Player, 0, 0, 3)
	_ = players.GiveWeapon(event.Player, int32(omp.WeaponM4), 500)
	_ = players.SetHealth(event.Player, 100)
	return nil
}

// handleHelpCommand shows available commands to the player.
func handleHelpCommand(player *omp.Player) bool {
	_ = players.SendClientMessage(player, uint32(omp.ColorYellow), "Available Commands:")
	_ = players.SendClientMessage(player, uint32(omp.ColorWhite), "/help - Show this help")
	_ = players.SendClientMessage(player, uint32(omp.ColorWhite), "/heal - Restore health")
	return true
}

// handleHealCommand restores player health.
func handleHealCommand(player *omp.Player) bool {
	_ = players.SetHealth(player, 100)
	_ = players.SendClientMessage(player, uint32(omp.ColorGreen), "Health restored!")
	return true
}

func (gm *FreeroamGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	command := event.Command.Clone()

	switch command {
	case "help":
		handleHelpCommand(event.Player)
	case "heal":
		handleHealCommand(event.Player)
	default:
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorRed), "Unknown command. Type /help for commands.")
	}

	return true, nil
}

func (gm *FreeroamGamemode) OnTick(ctx context.Context, event *omp.TickEvent) error {
	gm.tickCount++
	if gm.tickCount%60000 == 0 {
		minutes := gm.tickCount / 60000
		_ = core.Log(fmt.Sprintf("[Freeroam] Server uptime: %d minutes", minutes))
	}
	return nil
}

func NewGamemode() runtime.Gamemode {
	return &FreeroamGamemode{}
}

func init() {
	runtime.Bootstrap(
		runtime.WithComponentName("ompgo_freeroam"),
		runtime.WithGamemode(NewGamemode),
	)
}

func main() {}
