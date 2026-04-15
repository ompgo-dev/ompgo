// Package main implements a deathmatch gamemode for open.mp
// This gamemode features:
// - Random spawns around San Andreas
// - Automatic weapon loadout on spawn
// - Score tracking for kills
// - Death penalties
// - Player stats display
package main

import (
	"context"
	"fmt"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/all"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

// Deathmatch implements a deathmatch gamemode
type Deathmatch struct {
	omp.BaseEventHandler
	scores map[*omp.Player]int // Track player scores
}

// Spawn locations around San Andreas
var spawnLocations = []omp.Vector3{
	{X: 1958.3783, Y: 1343.1572, Z: 15.3746},  // Las Venturas
	{X: 2495.0857, Y: 1647.2380, Z: 10.8203},  // Las Venturas North
	{X: -1958.2280, Y: 128.6537, Z: 27.6875},  // San Fierro
	{X: -2188.7500, Y: 2405.8525, Z: 4.9688},  // San Fierro Bay
	{X: 1684.9200, Y: -2244.5679, Z: 13.5469}, // Los Santos Airport
	{X: 2240.1094, Y: -1258.8438, Z: 23.8516}, // Los Santos East
	{X: -2626.5500, Y: 1404.8700, Z: 7.1016},  // Bayside
	{X: 0.0, Y: 0.0, Z: 3.0},                  // Blueberry
}

// Weapon loadout: 9mm, Combat Shotgun, M4, Sniper, Grenades
var weaponLoadout = []struct {
	weaponID omp.WeaponID
	ammo     int32
}{
	{omp.WeaponColt45, 100},       // 9mm
	{omp.WeaponCombatShotgun, 50}, // Combat Shotgun
	{omp.WeaponM4, 200},           // M4
	{omp.WeaponSniper, 50},        // Sniper Rifle
	{omp.WeaponGrenade, 5},        // Grenades
}

// NewGamemode creates a new Deathmatch gamemode instance
func NewGamemode() runtime.Gamemode {
	return &Deathmatch{
		scores: make(map[*omp.Player]int),
	}
}

// OnLoad is called when the gamemode is loaded
func (dm *Deathmatch) OnLoad(ctx context.Context) error {
	fmt.Println("===========================================")
	fmt.Println("  Deathmatch Gamemode")
	fmt.Println("  Fight to be the best!")
	fmt.Println("===========================================")

	// Send welcome message to all players
	_ = all.SendClientMessage(
		uint32(omp.ColorGreen),
		"Deathmatch gamemode loaded! Fight to be the champion!",
	)

	fmt.Println("[Deathmatch] Gamemode loaded successfully")
	return nil
}

// OnPlayerConnect is called when a player connects
func (dm *Deathmatch) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}

	// Initialize player score
	dm.scores[event.Player] = 0

	// Get player name
	name, _ := players.GetName(event.Player)

	fmt.Printf("[Deathmatch] Player connected: %s\n", name)

	// Send welcome message
	_ = players.SendClientMessage(
		event.Player,
		uint32(omp.ColorYellow),
		fmt.Sprintf("Welcome %s! Type /help for commands", name),
	)
	_ = players.SendClientMessage(
		event.Player,
		uint32(omp.ColorWhite),
		"Kill other players to earn score. Good luck!",
	)
	return nil
}

// OnPlayerDisconnect is called when a player disconnects
func (dm *Deathmatch) OnPlayerDisconnect(ctx context.Context, event *omp.PlayerDisconnectEvent) error {
	if event.Player == nil {
		return nil
	}

	// Clean up player score
	delete(dm.scores, event.Player)

	fmt.Println("[Deathmatch] Player disconnected")
	return nil
}

// OnPlayerSpawn is called when a player spawns
func (dm *Deathmatch) OnPlayerSpawn(ctx context.Context, event *omp.PlayerSpawnEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}
	// Choose random spawn location
	spawn := spawnLocations[len(dm.scores)%len(spawnLocations)]

	// Set player position
	if !players.SetPos(event.Player, spawn.X, spawn.Y, spawn.Z) {
		fmt.Printf("[Deathmatch] Failed to set player spawn position\n")
	}

	// Set full health and armour
	_ = players.SetHealth(event.Player, 100.0)
	_ = players.SetArmor(event.Player, 100.0)

	// Give weapon loadout
	for _, weapon := range weaponLoadout {
		if !players.GiveWeapon(event.Player, int32(weapon.weaponID), weapon.ammo) {
			fmt.Printf("[Deathmatch] Failed to give weapon\n")
		}
	}

	// Send spawn message
	score := dm.scores[event.Player]
	_ = players.SendClientMessage(
		event.Player,
		uint32(omp.ColorGreen),
		fmt.Sprintf("You spawned! Score: %d | Kill enemies to earn points!", score),
	)

	fmt.Println("[Deathmatch] Player spawned with weapons")
	return nil
}

// OnPlayerDeath is called when a player dies
func (dm *Deathmatch) OnPlayerDeath(ctx context.Context, event *omp.PlayerDeathEvent) error {
	// Get killer and victim names
	var killerName, victimName string
	if event.Player != nil && event.Player.Valid() {
		victimName, _ = players.GetName(event.Player)
		// Death penalty
		if dm.scores[event.Player] > 0 {
			dm.scores[event.Player]--
		}
	}

	if event.Killer != nil && event.Killer.Valid() {
		killerName, _ = players.GetName(event.Killer)
		// Award kill point
		dm.scores[event.Killer]++

		// Notify killer
		score := dm.scores[event.Killer]
		_ = players.SendClientMessage(
			event.Killer,
			uint32(omp.ColorGreen),
			fmt.Sprintf("You killed %s! Score: %d (+1)", victimName, score),
		)
	}

	// Broadcast death message
	deathMsg := ""
	if event.Killer != nil && event.Killer.Valid() {
		deathMsg = fmt.Sprintf("%s killed %s", killerName, victimName)
	} else {
		deathMsg = fmt.Sprintf("%s died", victimName)
	}

	_ = all.SendClientMessage(uint32(omp.ColorGrey), deathMsg)
	fmt.Printf("[Deathmatch] %s\n", deathMsg)
	return nil
}

// OnPlayerText is called when a player sends a chat message
func (dm *Deathmatch) OnPlayerText(ctx context.Context, event *omp.PlayerTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	name, _ := players.GetName(event.Player)
	message := fmt.Sprintf("%s: %s", name, event.Text.Clone())

	_ = all.SendClientMessage(uint32(omp.ColorWhite), message)
	fmt.Printf("[Chat] %s\n", message)
	return false, nil
}

// OnPlayerCommandText is called when a player enters a command
func (dm *Deathmatch) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}
	cmd := event.Command.Clone()

	switch cmd {
	case "/help":
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "=== Deathmatch Commands ===")
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "/help - Show this help")
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "/score - Show your score")
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "/top - Show top players")
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "/kill - Suicide (death penalty applies)")
		return true, nil

	case "/score":
		score := dm.scores[event.Player]
		_ = players.SendClientMessage(
			event.Player,
			uint32(omp.ColorGreen),
			fmt.Sprintf("Your score: %d kills", score),
		)
		return true, nil

	case "/top":
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorYellow), "=== Top Players ===")
		// Sort and display top 5 (simplified - just showing concept)
		count := 0
		for player, score := range dm.scores {
			if count >= 5 {
				break
			}
			name, _ := players.GetName(player)
			_ = players.SendClientMessage(
				event.Player,
				uint32(omp.ColorWhite),
				fmt.Sprintf("%d. %s - %d kills", count+1, name, score),
			)
			count++
		}
		return true, nil

	case "/kill":
		// Set health to 0 to kill player
		_ = players.SetHealth(event.Player, 0.0)
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorRed), "You committed suicide!")
		return true, nil

	default:
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorRed), "Unknown command. Type /help for commands")
		return true, nil
	}
}
