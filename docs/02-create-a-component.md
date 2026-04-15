# Step 2: Create a component project

This step gets you to a runnable project layout.

## Recommended path: use the CLI scaffold

The fastest way to start is the `ompgo` CLI.

Create a new project:

```bash
go run github.com/ompgo-dev/ompgo/cmd/ompgo@latest init -name freeroam -module github.com/you/freeroam
```

That creates a folder named `freeroam` with:

- a `go.mod`
- a `main.go`
- a minimal `runtime.Bootstrap(...)` setup

If you want the project written somewhere else:

```bash
go run github.com/ompgo-dev/ompgo/cmd/ompgo@latest init -name freeroam -module github.com/you/freeroam -out /path/to/freeroam
```

If you prefer installing the CLI first:

```bash
go install github.com/ompgo-dev/ompgo/cmd/ompgo@latest
ompgo init -name freeroam -module github.com/you/freeroam
```

## Manual path

If you are adding ompgo to an existing module, create a `main.go` and add the dependency yourself.

Your module only needs a simple layout at first:

```text
your-component/
  go.mod
  main.go
```

Then add ompgo to the module:

```bash
go get github.com/ompgo-dev/ompgo@latest
```

## What to expect next

In the next step you will replace the scaffolded `main.go` with a minimal gamemode that:

- logs when the component loads
- welcomes a player when they connect
- registers itself through `runtime.Bootstrap(...)`

## Move on when

Move to the next step once you have a project directory with `go.mod` and `main.go`.

Next: [Step 3: Write your first gamemode](./03-first-gamemode.md)