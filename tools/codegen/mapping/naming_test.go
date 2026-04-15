package mapping

import "testing"

func TestSanitizeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "underscores", input: "player_id", want: "playerID"},
		{name: "ip address", input: "ipaddress", want: "ipAddress"},
		{name: "reserved type", input: "type", want: "typ"},
		{name: "reserved func", input: "func", want: "fn"},
		{name: "npc acronym", input: "npc", want: "npc"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := SanitizeName(tt.input); got != tt.want {
				t.Fatalf("SanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeGoName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "keyword not rewritten by sanitize", input: "var", want: "varValue"},
		{name: "keyword rewritten by sanitize", input: "type", want: "typ"},
		{name: "regular identifier", input: "player_name", want: "playerName"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := SafeGoName(tt.input); got != tt.want {
				t.Fatalf("SafeGoName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEventNamingHelpers(t *testing.T) {
	t.Parallel()

	if got := EventStructName("onPlayerConnect"); got != "PlayerConnectEvent" {
		t.Fatalf("EventStructName() = %q, want %q", got, "PlayerConnectEvent")
	}
	if got := EventMethodName("PlayerConnectEvent"); got != "OnPlayerConnect" {
		t.Fatalf("EventMethodName() = %q, want %q", got, "OnPlayerConnect")
	}
	if got := HandlerVarName("onPlayerConnect"); got != "onPlayerConnectHandlers" {
		t.Fatalf("HandlerVarName() = %q, want %q", got, "onPlayerConnectHandlers")
	}
	if got := RegisterName("onPlayerConnect"); got != "RegisterOnPlayerConnect" {
		t.Fatalf("RegisterName() = %q, want %q", got, "RegisterOnPlayerConnect")
	}
	if got := ExportName("onPlayerConnect"); got != "OMPGO_OnPlayerConnect" {
		t.Fatalf("ExportName() = %q, want %q", got, "OMPGO_OnPlayerConnect")
	}
}
