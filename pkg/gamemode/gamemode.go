// Package gamemode provides higher-level API helpers for gamemode code.
// Use pkg/omp for generated entities, events, constants, and the base event handler.
// Use pkg/runtime for bootstrap and the Gamemode interface.
//
// Usage Example:
//
//	api := gamemode.Initialize()
//	_ = api.Log("ready")
package gamemode
