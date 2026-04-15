# Step 5: Add commands and modules

This step shows how to grow beyond a single file without making the project hard to follow.

## Handle commands in the gamemode

Commands are handled in `OnPlayerCommandText`.

```go
func (gm *MyGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if event.Player == nil || !event.Player.Valid() {
		return false, nil
	}

	if event.Command.EqualString("/help") {
		_ = players.SendClientMessage(event.Player, uint32(omp.ColorWhite), "Available commands: /help, /heal")
		return true, nil
	}

	return false, nil
}
```

Use the boolean return value to tell open.mp whether the command was handled.

## Split behavior into modules

You do not need to put every handler on the main gamemode type. You can attach extra handler objects.

```go
type LoggingModule struct{}

func (m *LoggingModule) OnPlayerConnect(ctx context.Context, event *omp.PlayerConnectEvent) error {
	if event.Player == nil || !event.Player.Valid() {
		return nil
	}

	_ = players.SendClientMessage(event.Player, uint32(omp.ColorGreen), "Hello from module!")
	return nil
}

func (gm *MyGamemode) OnLoad(ctx context.Context) error {
	unregister := runtime.AttachHandlers(&LoggingModule{})
	_ = unregister
	return nil
}
```

This is useful when you want to keep features separate, such as:

- commands
- logging
- player onboarding
- admin tools

## Where this pattern fits

- Keep the main gamemode focused on startup and high-level flow.
- Put feature-specific handlers on smaller types.
- Use `pkg/omp/<group>` helpers from either place.

For a slightly larger example, read [examples/freeroam/main.go](../examples/freeroam/main.go).

## Move on when

Move on once you can handle a command and understand how to attach an extra handler object.

Next: [Step 6: Add context and error handling](./06-context-and-errors.md)