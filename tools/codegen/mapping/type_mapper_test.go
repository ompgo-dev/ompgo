package mapping

import "testing"

func TestResolveVoidPtrKind_EventOverrides(t *testing.T) {
	cases := []struct {
		event string
		name  string
		want  EntityKind
	}{
		{"onPlayerClickPlayer", "clicked", kindPlayer},
		{"onPlayerGiveDamage", "to", kindPlayer},
		{"onPlayerShotPlayer", "target", kindPlayer},
		{"onPlayerShotVehicle", "target", kindVehicle},
		{"onPlayerShotObject", "target", kindObject},
		{"onPlayerShotPlayerObject", "target", kindObject},
		{"onNPCTakeDamage", "damager", kindPlayer},
		{"onNPCGiveDamage", "damaged", kindPlayer},
	}

	for _, c := range cases {
		got, ok := ResolveVoidPtrKind(c.event, c.name)
		if !ok || got != c.want {
			t.Fatalf("resolveVoidPtrKind(%q, %q) = %v, %v; want %v, true", c.event, c.name, got, ok, c.want)
		}
	}
}

func TestResolveVoidPtrKind_ByName(t *testing.T) {
	cases := []struct {
		name string
		want EntityKind
	}{
		{"player", kindPlayer},
		{"killer", kindPlayer},
		{"vehicle", kindVehicle},
		{"trailer", kindVehicle},
		{"actor", kindActor},
		{"object", kindObject},
		{"playerobject", kindPlayerObject},
		{"pickup", kindPickup},
		{"textdraw", kindTextDraw},
		{"textlabel", kindTextLabel},
		{"gangzone", kindGangZone},
		{"menu", kindMenu},
		{"checkpoint", kindCheckpoint},
		{"class", kindClass},
		{"npc", kindNPC},
	}

	for _, c := range cases {
		got, ok := ResolveVoidPtrKind("", c.name)
		if !ok || got != c.want {
			t.Fatalf("resolveVoidPtrKind(\"\", %q) = %v, %v; want %v, true", c.name, got, ok, c.want)
		}
	}

	if _, ok := ResolveVoidPtrKind("", "unknown"); ok {
		t.Fatalf("resolveVoidPtrKind(\"\", \"unknown\") unexpectedly matched")
	}
}
