# Step 1: Understand the moving pieces

This step gives you the mental model for how an ompgo component is put together.

## What you are building

An ompgo project is a Go shared library that open.mp loads as a component. The runtime receives open.mp events, builds typed Go events from them, and calls your gamemode handlers.

The basic flow looks like this:

1. open.mp loads your compiled shared library.
2. `pkg/runtime` bootstraps the component.
3. The runtime receives C-API callbacks from open.mp.
4. Your gamemode methods handle those events in Go.

## Packages you will use most

- `pkg/omp`: events, entities, constants, and `omp.BaseEventHandler`
- `pkg/omp/<group>`: subsystem helpers like `players`, `core`, and `vehicles`
- `pkg/runtime`: `Bootstrap`, lifecycle hooks, error policies, and extra handler registration

Most new code should import from `pkg/omp` and `pkg/runtime`.

## What a minimal component needs

You need four things:

1. A `main` package.
2. A gamemode type that embeds `omp.BaseEventHandler`.
3. One or more event methods such as `OnLoad` or `OnPlayerConnect`.
4. A call to `runtime.Bootstrap(...)` in `init()`.

## Prerequisites

Before you start, make sure you have:

- Go 1.25+
- CGO enabled with `CGO_ENABLED=1`
- An open.mp server with the C-API component available as `$CAPI.so` or `$CAPI.dll`

## How to use this guide

- If you want the shortest route to a working project, continue to [Step 2](./02-create-a-component.md).
- If you already created a project and only need code structure, skip to [Step 3](./03-first-gamemode.md).

## Move on when

Move to the next step once you understand that `pkg/runtime` handles the component wiring and your code mostly lives in gamemode methods.

Next: [Step 2: Create a component project](./02-create-a-component.md)