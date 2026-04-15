# ompgo docs

This folder is organized as a step-by-step learning path.

If you are new to ompgo, follow the steps in order:

1. [Step 1: Understand the moving pieces](./01-overview.md)
2. [Step 2: Create a component project](./02-create-a-component.md)
3. [Step 3: Write your first gamemode](./03-first-gamemode.md)
4. [Step 4: Build and install it](./04-build-and-install.md)
5. [Step 5: Add commands and modules](./05-commands-and-modules.md)
6. [Step 6: Add context and error handling](./06-context-and-errors.md)

## Choose your path

- Want the fastest way to a working component: start at Step 2.
- Want to understand the package layout before writing code: start at Step 1.
- Already have a working component and want runtime diagnostics: jump to Step 6.

## Examples path

If you want to learn by reading code alongside the docs, use [examples/README.md](../examples/README.md).

The recommended order is:

1. [examples/basic/main.go](../examples/basic/main.go)
2. [examples/events/main.go](../examples/events/main.go)
3. [examples/freeroam/main.go](../examples/freeroam/main.go)
4. [examples/context/main.go](../examples/context/main.go)
5. [examples/deathmatch/main.go](../examples/deathmatch/main.go)
6. [examples/grandlarc/main.go](../examples/grandlarc/main.go)

## Quick reference

- `pkg/omp`: generated events, entities, constants, and `omp.BaseEventHandler`
- `pkg/omp/<group>`: helper functions such as `players`, `core`, `vehicles`, or `textdraw`
- `pkg/runtime`: bootstrap, lifecycle hooks, handler attachment, and runtime policies
- `pkg/gamemode`: optional higher-level compatibility layer