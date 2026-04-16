ompgo
======

Go runtime and bindings for building open.mp components with CGO.

What this is
- Go runtime that loads open.mp C-API and dispatches events to your gamemode.
- Generated Go bindings for open.mp functions and events.
- A simple pattern for building a single Go shared library component.

Quick start
1. Write a `main` package that registers a gamemode and exports the required component entry points.
2. Build with `-buildmode=c-shared`.
3. Copy the .so into your server’s components folder.

Docs
- [docs/README.md](docs/README.md)
- [docs/using-ompgo.md](docs/using-ompgo.md)

Dependency workflow
- The upstream open.mp C API now lives in the git submodule at `third_party/openmp-capi`.
- The repo builds from committed snapshots in `tools/codegen/data/openmp-capi` and `pkg/runtime/include/ompcapi.h`.
- Use `task sync:openmp-capi` to copy `api.json`, `events.json`, and `ompcapi.h` from the submodule into those committed paths.
- Use `task update:openmp-capi` to update the submodule and refresh the committed snapshots in one step.

Examples
- [examples/README.md](examples/README.md)
- [examples/basic/main.go](examples/basic/main.go)
- [examples/context/main.go](examples/context/main.go)
- [examples/deathmatch/main.go](examples/deathmatch/main.go)
- [examples/events/main.go](examples/events/main.go)
- [examples/freeroam/main.go](examples/freeroam/main.go)
- [examples/grandlarc/main.go](examples/grandlarc/main.go)

Packages
  - [ompgo-streamer](https://github.com/ompgo-dev/ompgo-streamer) - Recreates the streamer plugin using ompgo
  - [ompgo-extras](https://github.com/ompgo-dev/ompgo-extras) - A bunch of helper packages for common gamemode workflows. (command and dialog routers, and textdraw builder)
