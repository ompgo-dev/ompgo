// Package main implements the classic Grand Larceny gamemode for open.mp
package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/all"
	"github.com/ompgo-dev/ompgo/pkg/omp/class"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/omp/textdraw"
	"github.com/ompgo-dev/ompgo/pkg/omp/vehicle"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

import "C"

type city int

const (
	cityLosSantos city = iota
	citySanFierro
	cityLasVenturas
)

type playerData struct {
	selectedCity      city
	selectedCitySet   bool
	hasCitySelected   bool
	lastSelectionTick time.Time
}

type GrandLarc struct {
	omp.BaseEventHandler
	playersData map[int32]*playerData
	spawnLocs   spawnLocations
	rng         *rand.Rand

	classSelectionHelper *omp.TextDraw
	losSantosTD          *omp.TextDraw
	sanFierroTD          *omp.TextDraw
	lasVenturasTD        *omp.TextDraw
}

func NewGamemode() runtime.Gamemode {
	return &GrandLarc{
		playersData: make(map[int32]*playerData),
		spawnLocs:   newSpawnLocations(),
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func init() {
	runtime.Bootstrap(
		runtime.WithComponentName("Grand Larceny"),
		runtime.WithComponentVersion(runtime.Version{
			Major:  1,
			Minor:  0,
			Patch:  0,
			Prerel: 0,
		}),
		runtime.WithGamemode(NewGamemode),
	)
}

func (g *GrandLarc) OnLoad(ctx context.Context) error {
	_ = core.GameModeSetText("Grand Larceny")
	_ = core.ShowPlayerMarkers(1)
	_ = core.ShowNameTags(true)
	_ = core.SetNameTagsDrawDistance(40.0)
	_ = all.EnableStuntBonus(false)
	_ = core.DisableEntryExitMarkers()
	_ = core.SetWeather(2)
	_ = core.SetWorldTime(11)

	g.classSelectionHelper = createHelperTextDraw()
	g.losSantosTD = createCityNameTextDraw("Los Santos")
	g.sanFierroTD = createCityNameTextDraw("San Fierro")
	g.lasVenturasTD = createCityNameTextDraw("Las Venturas")

	createAllClasses()
	g.loadStaticVehicles()
	return nil
}

func (g *GrandLarc) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}

	_ = players.ShowGameText(event.Player, "~w~Grand Larceny", 3000, 4)
	_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "Welcome to {88AA88}G{FFFFFF}rand {88AA88}L{FFFFFF}arceny")

	playerID := players.GetID(event.Player)
	g.playersData[playerID] = &playerData{
		selectedCitySet:   false,
		hasCitySelected:   false,
		lastSelectionTick: time.Now(),
	}
	return nil
}

func (g *GrandLarc) OnPlayerDisconnect(ctx context.Context, event *omp.PlayerDisconnectEvent) error {
	if event.Player == nil {
		return nil
	}
	delete(g.playersData, players.GetID(event.Player))
	return nil
}

func (g *GrandLarc) OnPlayerSpawn(ctx context.Context, event *omp.PlayerSpawnEvent) error {
	if event.Player == nil || !event.Player.Valid() || players.IsNPC(event.Player) {
		return nil
	}

	_ = players.SetInterior(event.Player, 0)
	_ = players.ToggleClock(event.Player, false)
	_ = players.ResetMoney(event.Player)
	_ = players.GiveMoney(event.Player, 30000)

	pdata := g.ensurePlayerData(event.Player)
	if !pdata.selectedCitySet {
		pdata.selectedCity = cityLosSantos
		pdata.selectedCitySet = true
	}

	switch pdata.selectedCity {
	case cityLosSantos:
		spawn := g.spawnLocs.randomLS(g.rng)
		_ = players.SetPos(event.Player, spawn.x, spawn.y, spawn.z)
		_ = players.SetFacingAngle(event.Player, spawn.angle)
	case citySanFierro:
		spawn := g.spawnLocs.randomSF(g.rng)
		_ = players.SetPos(event.Player, spawn.x, spawn.y, spawn.z)
		_ = players.SetFacingAngle(event.Player, spawn.angle)
	case cityLasVenturas:
		spawn := g.spawnLocs.randomLV(g.rng)
		_ = players.SetPos(event.Player, spawn.x, spawn.y, spawn.z)
		_ = players.SetFacingAngle(event.Player, spawn.angle)
	}

	_ = players.GiveWeapon(event.Player, int32(omp.WeaponColt45), 100)
	_ = players.ToggleClock(event.Player, false)
	return nil
}

func (g *GrandLarc) OnPlayerDeath(ctx context.Context, event *omp.PlayerDeathEvent) error {
	if event.Player == nil {
		return nil
	}

	pdata := g.ensurePlayerData(event.Player)
	pdata.hasCitySelected = false

	if event.Killer != nil && event.Killer.Valid() {
		playerCash := players.GetMoney(event.Player)
		if playerCash > 0 {
			_ = players.GiveMoney(event.Killer, playerCash)
			_ = players.ResetMoney(event.Player)
		}
	} else {
		_ = players.ResetMoney(event.Player)
	}
	return nil
}

func (g *GrandLarc) OnPlayerRequestClass(ctx context.Context, event *omp.PlayerRequestClassEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() || players.IsNPC(event.Player) {
		return true, nil
	}

	pdata := g.ensurePlayerData(event.Player)
	if pdata.hasCitySelected {
		g.setupCharSelection(event.Player)
		return true, nil
	}

	if players.GetState(event.Player) != int32(omp.PlayerStateSpectating) {
		_ = players.ToggleSpectating(event.Player, true)
		g.classSelectionHelperShow(event.Player)
		pdata.selectedCitySet = false
	}

	return false, nil
}

func (g *GrandLarc) OnPlayerUpdate(ctx context.Context, event *omp.PlayerUpdateEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() || players.IsNPC(event.Player) {
		return true, nil
	}

	pdata := g.ensurePlayerData(event.Player)
	if !pdata.hasCitySelected && players.GetState(event.Player) == int32(omp.PlayerStateSpectating) {
		g.handleCitySelection(event.Player, pdata)
		return true, nil
	}

	if players.GetWeapon(event.Player) == int32(omp.WeaponMinigun) {
		_ = players.Kick(event.Player)
		return false, nil
	}

	return true, nil
}

func (g *GrandLarc) ensurePlayerData(player *omp.Player) *playerData {
	playerID := players.GetID(player)
	pdata, ok := g.playersData[playerID]
	if ok {
		return pdata
	}

	pdata = &playerData{lastSelectionTick: time.Now()}
	g.playersData[playerID] = pdata
	return pdata
}

func (g *GrandLarc) setupCharSelection(player *omp.Player) {
	pdata := g.ensurePlayerData(player)
	if !pdata.selectedCitySet {
		return
	}

	switch pdata.selectedCity {
	case cityLosSantos:
		_ = players.SetInterior(player, 11)
		_ = players.SetPos(player, 508.7362, -87.4335, 998.9609)
		_ = players.SetFacingAngle(player, 0.0)
		_ = players.SetCameraPos(player, 508.7362, -83.4335, 998.9609)
		_ = players.SetCameraLookAt(player, 508.7362, -87.4335, 998.9609, int32(omp.CameraMove))
	case citySanFierro:
		_ = players.SetInterior(player, 3)
		_ = players.SetPos(player, -2673.8381, 1399.7424, 918.3516)
		_ = players.SetFacingAngle(player, 181.0)
		_ = players.SetCameraPos(player, -2673.2776, 1394.3859, 918.3516)
		_ = players.SetCameraLookAt(player, -2673.8381, 1399.7424, 918.3516, int32(omp.CameraMove))
	case cityLasVenturas:
		_ = players.SetInterior(player, 3)
		_ = players.SetPos(player, 349.0453, 193.2271, 1014.1797)
		_ = players.SetFacingAngle(player, 286.25)
		_ = players.SetCameraPos(player, 352.9164, 194.5702, 1014.1875)
		_ = players.SetCameraLookAt(player, 349.0453, 193.2271, 1014.1797, int32(omp.CameraMove))
	}
}

func (g *GrandLarc) setupSelectedCity(player *omp.Player) {
	pdata := g.ensurePlayerData(player)
	if !pdata.selectedCitySet {
		pdata.selectedCity = cityLosSantos
		pdata.selectedCitySet = true
	}

	switch pdata.selectedCity {
	case cityLosSantos:
		_ = players.SetInterior(player, 0)
		_ = players.SetCameraPos(player, 1630.6136, -2286.0298, 110.0)
		_ = players.SetCameraLookAt(player, 1887.6034, -1682.1442, 47.6167, int32(omp.CameraMove))
		g.showCityTextDraw(player, g.losSantosTD)
		g.hideCityTextDraw(player, g.sanFierroTD)
		g.hideCityTextDraw(player, g.lasVenturasTD)
	case citySanFierro:
		_ = players.SetInterior(player, 0)
		_ = players.SetCameraPos(player, -1300.8754, 68.0546, 129.4823)
		_ = players.SetCameraLookAt(player, -1817.9412, 769.3878, 132.6589, int32(omp.CameraMove))
		g.hideCityTextDraw(player, g.losSantosTD)
		g.showCityTextDraw(player, g.sanFierroTD)
		g.hideCityTextDraw(player, g.lasVenturasTD)
	case cityLasVenturas:
		_ = players.SetInterior(player, 0)
		_ = players.SetCameraPos(player, 1310.6155, 1675.9182, 110.739)
		_ = players.SetCameraLookAt(player, 2285.2944, 1919.3756, 68.2275, int32(omp.CameraMove))
		g.hideCityTextDraw(player, g.losSantosTD)
		g.hideCityTextDraw(player, g.sanFierroTD)
		g.showCityTextDraw(player, g.lasVenturasTD)
	}
}

func (g *GrandLarc) switchToNextCity(player *omp.Player, pdata *playerData) {
	if !pdata.selectedCitySet {
		pdata.selectedCity = cityLosSantos
		pdata.selectedCitySet = true
	} else {
		switch pdata.selectedCity {
		case cityLosSantos:
			pdata.selectedCity = citySanFierro
		case citySanFierro:
			pdata.selectedCity = cityLasVenturas
		case cityLasVenturas:
			pdata.selectedCity = cityLosSantos
		}
	}

	_ = players.PlayGameSound(player, 1052, 0, 0, 0)
	pdata.lastSelectionTick = time.Now()
	g.setupSelectedCity(player)
}

func (g *GrandLarc) switchToPreviousCity(player *omp.Player, pdata *playerData) {
	if !pdata.selectedCitySet {
		return
	}

	switch pdata.selectedCity {
	case cityLosSantos:
		pdata.selectedCity = cityLasVenturas
	case citySanFierro:
		pdata.selectedCity = cityLosSantos
	case cityLasVenturas:
		pdata.selectedCity = citySanFierro
	}

	_ = players.PlayGameSound(player, 1053, 0, 0, 0)
	pdata.lastSelectionTick = time.Now()
	g.setupSelectedCity(player)
}

func (g *GrandLarc) handleCitySelection(player *omp.Player, pdata *playerData) {
	if !pdata.selectedCitySet {
		g.switchToNextCity(player, pdata)
		return
	}

	if time.Since(pdata.lastSelectionTick) < 500*time.Millisecond {
		return
	}

	var keys int32
	var updown int32
	var leftright int32
	_ = players.GetKeys(player, &keys, &updown, &leftright)

	if keys&int32(omp.KeyFire) != 0 {
		pdata.hasCitySelected = true
		g.hideCityTextDraw(player, g.losSantosTD)
		g.hideCityTextDraw(player, g.sanFierroTD)
		g.hideCityTextDraw(player, g.lasVenturasTD)
		g.classSelectionHelperHide(player)
		_ = players.ToggleSpectating(player, false)
		return
	}

	if leftright > 0 {
		g.switchToNextCity(player, pdata)
		return
	}
	if leftright < 0 {
		g.switchToPreviousCity(player, pdata)
	}
}

func (g *GrandLarc) classSelectionHelperShow(player *omp.Player) {
	if g.classSelectionHelper != nil {
		_ = textdraw.ShowForPlayer(player, g.classSelectionHelper)
	}
}

func (g *GrandLarc) classSelectionHelperHide(player *omp.Player) {
	if g.classSelectionHelper != nil {
		_ = textdraw.HideForPlayer(player, g.classSelectionHelper)
	}
}

func (g *GrandLarc) showCityTextDraw(player *omp.Player, td *omp.TextDraw) {
	if td != nil {
		_ = textdraw.ShowForPlayer(player, td)
	}
}

func (g *GrandLarc) hideCityTextDraw(player *omp.Player, td *omp.TextDraw) {
	if td != nil {
		_ = textdraw.HideForPlayer(player, td)
	}
}

func createCityNameTextDraw(text string) *omp.TextDraw {
	var id int32
	td := textdraw.Create(10.0, 380.0, text, &id)
	if td == nil {
		return nil
	}
	_ = textdraw.SetUseBox(td, false)
	_ = textdraw.SetLetterSize(td, 1.25, 3.0)
	_ = textdraw.SetFont(td, int32(omp.TextDrawFontBeckettRegular))
	_ = textdraw.SetShadow(td, 0)
	_ = textdraw.SetOutline(td, 1)
	_ = textdraw.SetColor(td, 0xEEEEEEFF)
	return td
}

func createHelperTextDraw() *omp.TextDraw {
	var id int32
	td := textdraw.Create(10.0, 415.0, " Press ~b~~k~~GO_LEFT~ ~w~or ~b~~k~~GO_RIGHT~ ~w~to switch cities.~n~ Press ~r~~k~~PED_FIREWEAPON~ ~w~to select.", &id)
	if td == nil {
		return nil
	}
	_ = textdraw.SetUseBox(td, true)
	_ = textdraw.SetBoxColor(td, 0x222222BB)
	_ = textdraw.SetLetterSize(td, 0.3, 1.0)
	_ = textdraw.SetTextSize(td, 400.0, 40.0)
	_ = textdraw.SetFont(td, int32(omp.TextDrawFontBankGothic))
	_ = textdraw.SetShadow(td, 0)
	_ = textdraw.SetOutline(td, 1)
	_ = textdraw.SetBackgroundColor(td, 0x000000FF)
	_ = textdraw.SetColor(td, 0xFFFFFFFF)
	return td
}

func createAllClasses() {
	skins := []int32{
		298, 299, 300, 301, 302, 303, 304, 305, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289,
		265, 266, 267, 268, 269, 270, 1, 2, 3, 4, 5, 6, 8, 42, 65, 86, 119, 149, 208, 273, 289, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 68, 69, 70, 71, 72, 73, 75, 76, 78, 79, 80, 81,
		82, 83, 84, 85, 87, 88, 89, 91, 92, 93, 95, 96, 97, 98, 99,
	}

	for _, skin := range skins {
		var id int32
		_ = class.Add(255, skin, 1759.0189, -1898.1260, 13.5622, 266.4503, 0, 0, 0, 0, 0, 0, &id)
	}
}

func (g *GrandLarc) loadStaticVehicles() {
	vehicleFiles := []string{
		"trains",
		"pilots",
		"lv_law",
		"lv_airport",
		"lv_gen",
		"sf_law",
		"sf_airport",
		"sf_gen",
		"ls_law",
		"ls_airport",
		"ls_gen_inner",
		"ls_gen_outer",
		"whetstone",
		"bone",
		"flint",
		"tierra",
		"red_county",
	}

	total := 0
	for _, name := range vehicleFiles {
		path := filepath.Join("scriptfiles", "vehicles", name+".txt")
		count, err := loadStaticVehiclesFromFile(path)
		if err != nil {
			_ = core.Log(fmt.Sprintf("[GrandLarc] vehicle file %s: %v", path, err))
			continue
		}
		total += count
	}

	_ = core.Log(fmt.Sprintf("[GrandLarc] Total vehicles from files: %d", total))
}

func loadStaticVehiclesFromFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 7 {
			continue
		}

		modelid, err := parseInt32(parts[0])
		if err != nil {
			continue
		}
		x, err := parseFloat32(parts[1])
		if err != nil {
			continue
		}
		y, err := parseFloat32(parts[2])
		if err != nil {
			continue
		}
		z, err := parseFloat32(parts[3])
		if err != nil {
			continue
		}
		rotation, err := parseFloat32(parts[4])
		if err != nil {
			continue
		}
		colour1, err := parseInt32(parts[5])
		if err != nil {
			continue
		}
		colour2Field := strings.Fields(parts[6])
		if len(colour2Field) == 0 {
			continue
		}
		colour2, err := parseInt32(colour2Field[0])
		if err != nil {
			continue
		}

		var id int32
		_ = vehicle.Create(modelid, x, y, z, rotation, colour1, colour2, 30*60, false, &id)
		count++
	}

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}

func parseInt32(value string) (int32, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(value), 10, 32)
	return int32(v), err
}

func parseFloat32(value string) (float32, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(value), 32)
	return float32(v), err
}

func main() {}
