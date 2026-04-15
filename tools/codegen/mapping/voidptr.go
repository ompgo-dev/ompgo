package mapping

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

type EntityKind int

const (
	kindUnknown EntityKind = iota
	kindPlayer
	kindVehicle
	kindActor
	kindObject
	kindPlayerObject
	kindPickup
	kindTextDraw
	kindPlayerTextDraw
	kindTextLabel
	kindPlayerTextLabel
	kindGangZone
	kindMenu
	kindCheckpoint
	kindClass
	kindNPC
)

var eventVoidPtrOverrides = map[string]map[string]EntityKind{
	"onPlayerClickPlayer": {
		"clicked": kindPlayer,
	},
	"onPlayerGiveDamage": {
		"to": kindPlayer,
	},
	"onPlayerShotPlayer": {
		"target": kindPlayer,
	},
	"onPlayerShotVehicle": {
		"target": kindVehicle,
	},
	"onPlayerShotObject": {
		"target": kindObject,
	},
	"onPlayerShotPlayerObject": {
		"target": kindObject,
	},
	"onNPCTakeDamage": {
		"damager": kindPlayer,
	},
	"onNPCGiveDamage": {
		"damaged": kindPlayer,
	},
}

var apiVoidPtrOverrides = map[string]map[string]EntityKind{
	"player_senddeathmessage": {
		"killee": kindPlayer,
	},
	"all_senddeathmessage": {
		"killee": kindPlayer,
	},
	"player_sendmessagetoplayer": {
		"sender": kindPlayer,
	},
	"player_playcrimereport": {
		"suspect": kindPlayer,
	},
	"player_isstreamedin": {
		"other": kindPlayer,
	},
	"player_shownametagforplayer": {
		"other": kindPlayer,
	},
	"player_setmarkerforplayer": {
		"other": kindPlayer,
	},
	"player_getmarkerforplayer": {
		"other": kindPlayer,
	},
	"player_spectateplayer": {
		"target": kindPlayer,
	},
	"player_spectatevehicle": {
		"target": kindVehicle,
	},
	"playerobject_attachtoobject": {
		"attachedto": kindPlayerObject,
	},
	"object_attachtoobject": {
		"objattachedto": kindObject,
	},
}

var playerExactNames = map[string]struct{}{
	"killer":  {},
	"from":    {},
	"issuer":  {},
	"killee":  {},
	"sender":  {},
	"suspect": {},
	"other":   {},
}

func ResolveVoidPtrKind(eventName, paramName string) (EntityKind, bool) {
	paramLower := strings.ToLower(paramName)
	if eventName != "" {
		if overrides, ok := eventVoidPtrOverrides[eventName]; ok {
			if kind, ok := overrides[paramLower]; ok {
				return kind, true
			}
		}
	}

	return ResolveVoidPtrKindByName(paramLower)
}

func ResolveVoidPtrKindForAPI(fnName, paramName string) (EntityKind, bool) {
	nameLower := strings.ToLower(strings.TrimSpace(fnName))
	paramLower := strings.ToLower(paramName)
	if strings.HasPrefix(nameLower, "playerobject_") && strings.Contains(paramLower, "object") {
		return kindPlayerObject, true
	}
	if strings.HasPrefix(nameLower, "playertextdraw_") && strings.Contains(paramLower, "textdraw") {
		return kindPlayerTextDraw, true
	}
	if strings.HasPrefix(nameLower, "playertextlabel_") && strings.Contains(paramLower, "textlabel") {
		return kindPlayerTextLabel, true
	}
	if overrides, ok := apiVoidPtrOverrides[nameLower]; ok {
		if kind, ok := overrides[paramLower]; ok {
			return kind, true
		}
	}
	return ResolveVoidPtrKindByName(paramLower)
}

func ResolveVoidPtrKindByName(paramLower string) (EntityKind, bool) {
	if paramLower == "cls" {
		return kindClass, true
	}
	if _, ok := playerExactNames[paramLower]; ok {
		return kindPlayer, true
	}
	if strings.Contains(paramLower, "trailer") {
		return kindVehicle, true
	}
	if strings.Contains(paramLower, "playerobject") {
		return kindPlayerObject, true
	}
	if strings.Contains(paramLower, "playertextdraw") {
		return kindPlayerTextDraw, true
	}
	if strings.Contains(paramLower, "textdraw") || paramLower == "td" {
		return kindTextDraw, true
	}
	if strings.Contains(paramLower, "playertextlabel") {
		return kindPlayerTextLabel, true
	}
	if strings.Contains(paramLower, "textlabel") || strings.Contains(paramLower, "label") {
		return kindTextLabel, true
	}
	if strings.Contains(paramLower, "menu") {
		return kindMenu, true
	}
	if strings.Contains(paramLower, "checkpoint") {
		return kindCheckpoint, true
	}
	if strings.Contains(paramLower, "gangzone") || strings.Contains(paramLower, "zone") {
		return kindGangZone, true
	}
	if strings.Contains(paramLower, "vehicle") {
		return kindVehicle, true
	}
	if strings.Contains(paramLower, "actor") {
		return kindActor, true
	}
	if strings.Contains(paramLower, "object") {
		return kindObject, true
	}
	if strings.Contains(paramLower, "pickup") {
		return kindPickup, true
	}
	if strings.Contains(paramLower, "class") {
		return kindClass, true
	}
	if strings.Contains(paramLower, "npc") {
		return kindNPC, true
	}
	if strings.Contains(paramLower, "player") {
		return kindPlayer, true
	}

	return kindUnknown, false
}

func KindToHelperType(kind EntityKind) string {
	switch kind {
	case kindPlayer:
		return "*omp.Player"
	case kindVehicle:
		return "*omp.Vehicle"
	case kindActor:
		return "*omp.Actor"
	case kindObject:
		return "*omp.Object"
	case kindPlayerObject:
		return "*omp.PlayerObject"
	case kindPickup:
		return "*omp.Pickup"
	case kindTextDraw:
		return "*omp.TextDraw"
	case kindPlayerTextDraw:
		return "*omp.PlayerTextDraw"
	case kindTextLabel:
		return "*omp.TextLabel"
	case kindPlayerTextLabel:
		return "*omp.PlayerTextLabel"
	case kindGangZone:
		return "*omp.GangZone"
	case kindMenu:
		return "*omp.Menu"
	case kindCheckpoint:
		return "*omp.Checkpoint"
	case kindClass:
		return "*omp.Class"
	case kindNPC:
		return "*omp.NPC"
	default:
		return "unsafe.Pointer"
	}
}

func KindToEventGoType(kind EntityKind) jen.Code {
	switch kind {
	case kindPlayer:
		return jen.Op("*").Id("Player")
	case kindVehicle:
		return jen.Op("*").Id("Vehicle")
	case kindActor:
		return jen.Op("*").Id("Actor")
	case kindObject:
		return jen.Op("*").Id("Object")
	case kindPlayerObject:
		return jen.Op("*").Id("PlayerObject")
	case kindPickup:
		return jen.Op("*").Id("Pickup")
	case kindTextDraw:
		return jen.Op("*").Id("TextDraw")
	case kindPlayerTextDraw:
		return jen.Op("*").Id("PlayerTextDraw")
	case kindTextLabel:
		return jen.Op("*").Id("TextLabel")
	case kindPlayerTextLabel:
		return jen.Op("*").Id("PlayerTextLabel")
	case kindGangZone:
		return jen.Op("*").Id("GangZone")
	case kindMenu:
		return jen.Op("*").Id("Menu")
	case kindCheckpoint:
		return jen.Op("*").Id("Checkpoint")
	case kindClass:
		return jen.Op("*").Id("Class")
	case kindNPC:
		return jen.Op("*").Id("NPC")
	default:
		return jen.Interface()
	}
}

func KindToConstructor(kind EntityKind) (string, bool) {
	switch kind {
	case kindPlayer:
		return "omp.NewPlayer", true
	case kindVehicle:
		return "omp.NewVehicle", true
	case kindActor:
		return "omp.NewActor", true
	case kindObject:
		return "omp.NewObject", true
	case kindPlayerObject:
		return "omp.NewPlayerObject", true
	case kindPickup:
		return "omp.NewPickup", true
	case kindTextDraw:
		return "omp.NewTextDraw", true
	case kindPlayerTextDraw:
		return "omp.NewPlayerTextDraw", true
	case kindTextLabel:
		return "omp.NewTextLabel", true
	case kindPlayerTextLabel:
		return "omp.NewPlayerTextLabel", true
	case kindGangZone:
		return "omp.NewGangZone", true
	case kindMenu:
		return "omp.NewMenu", true
	case kindCheckpoint:
		return "omp.NewCheckpoint", true
	case kindClass:
		return "omp.NewClass", true
	case kindNPC:
		return "omp.NewNPC", true
	default:
		return "", false
	}
}

func ConstructorForVoidPtr(eventName, paramName string) (string, bool) {
	kind, ok := ResolveVoidPtrKind(eventName, paramName)
	if !ok {
		return "", false
	}
	return KindToConstructor(kind)
}
