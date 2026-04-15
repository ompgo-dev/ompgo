# Step 3: Write your first gamemode

This step gives you the smallest useful ompgo component.

## Start with this code

```go
package main

import (
	"context"

	"github.com/ompgo-dev/ompgo/pkg/omp"
	"github.com/ompgo-dev/ompgo/pkg/omp/core"
	"github.com/ompgo-dev/ompgo/pkg/omp/players"
	"github.com/ompgo-dev/ompgo/pkg/runtime"
)

type MyGamemode struct {
	omp.BaseEventHandler
}

func (gm *MyGamemode) OnLoad(ctx context.Context) error {
	_ = core.Log("[MyGamemode] Loaded")
	return nil
}

func (gm *MyGamemode) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}

	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "Welcome!")
	return nil
}

func NewGamemode() runtime.Gamemode {
	return &MyGamemode{}
}

func init() {
	runtime.Bootstrap(
		runtime.WithComponentName("ompgo_mygamemode"),
		runtime.WithGamemode(NewGamemode),
	)
}

func main() {}
```

## What each part does

`omp.BaseEventHandler`

Gives you default no-op implementations for all events, so you only implement the ones you care about.

`OnLoad`

Runs when the gamemode is created. This is a good place for startup logging and in-memory setup.

`OnPlayerConnect`

Shows a basic gameplay event handler. Always validate event-owned objects like `event.Player` before using them.

`runtime.Bootstrap(...)`

Registers the component and wires the runtime. You do not need to hand-write `Init`, `ComponentEntryPoint`, or `ComponentCleanup` exports.

## Rules to keep in mind

- Event handlers use `context.Context`.
- Most handlers return `error`.
- Blocking events return `(bool, error)`.
- Use `pkg/omp/<group>` helper packages for server functions.

## Compare against the example

If you want to see a slightly fuller version of this pattern, read [examples/basic/main.go](../examples/basic/main.go).

## Move on when

Move on once your `main.go` has a gamemode type, at least one event handler, and a `runtime.Bootstrap(...)` call.

Next: [Step 4: Build and install it](./04-build-and-install.md)