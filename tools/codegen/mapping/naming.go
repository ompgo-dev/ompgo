package mapping

import "strings"

// Capitalize uppercases the first letter and normalizes common abbreviations.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	switch strings.ToLower(s) {
	case "ip":
		return "IP"
	case "id":
		return "ID"
	case "npc":
		return "NPC"
	case "rcon":
		return "RCON"
	case "td":
		return "TD"
	case "vw":
		return "VW"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// IsGoKeyword reports whether word is a Go keyword.
func IsGoKeyword(word string) bool {
	switch word {
	case "break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var":
		return true
	default:
		return false
	}
}

// SanitizeName normalizes identifiers to Go-friendly camelCase and fixes common acronyms.
func SanitizeName(name string) string {
	// Special case: ipaddress should become ipAddress (then capitalize handles IP)
	if strings.ToLower(name) == "ipaddress" {
		return "ipAddress"
	}

	// Convert underscores to camelCase
	if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		for i := 1; i < len(parts); i++ {
			parts[i] = Capitalize(parts[i])
		}
		name = strings.Join(parts, "")
	}

	// Fix common Go naming conventions
	name = strings.ReplaceAll(name, "Id", "ID")
	name = strings.ReplaceAll(name, "Ip", "IP")

	// Handle reserved keywords
	switch name {
	case "type":
		return "typ"
	case "range":
		return "rng"
	case "interface":
		return "iface"
	case "func":
		return "fn"
	default:
		return name
	}
}

// SafeGoName sanitizes and avoids Go keywords for identifiers.
func SafeGoName(name string) string {
	clean := SanitizeName(name)
	if IsGoKeyword(clean) {
		return clean + "Value"
	}
	return clean
}

// EventStructName builds the Go struct name for an event.
func EventStructName(eventName string) string {
	name := strings.TrimPrefix(eventName, "on")
	return Capitalize(name) + "Event"
}

// EventMethodName builds the handler method name from a struct name.
func EventMethodName(structName string) string {
	return strings.TrimSuffix("On"+structName, "Event")
}

// HandlerName converts event name to handler method name.
func HandlerName(eventName string) string {
	name := strings.TrimPrefix(eventName, "on")
	return "On" + name
}

// HandlerVarName derives the handler variable name from event name.
func HandlerVarName(eventName string) string {
	name := HandlerName(eventName)
	if name == "" {
		return "handlers"
	}
	return strings.ToLower(name[:1]) + name[1:] + "Handlers"
}

// RegisterName builds the register function name for an event.
func RegisterName(eventName string) string {
	return "Register" + HandlerName(eventName)
}

// ExportName gets the C export function name.
// onPlayerConnect -> OMPGO_OnPlayerConnect
func ExportName(eventName string) string {
	name := strings.TrimPrefix(eventName, "on")
	return "OMPGO_On" + name
}
