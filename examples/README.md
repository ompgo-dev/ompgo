# ompgo examples

This folder is organized as a learning path, not just a dump of sample projects.

If you want to learn by reading code, go through these examples in order.

## Step 1: Start from the smallest working example

- [basic/main.go](./basic/main.go)

This is the best first example to read. It shows:

- a minimal gamemode type
- `runtime.Bootstrap(...)`
- a few common event handlers
- simple player messaging

## Step 2: See a focused event example

- [events/main.go](./events/main.go)

Read this after `basic` if you want a smaller event-focused sample that stays close to the core runtime shape.

## Step 3: Read a slightly larger gameplay flow

- [freeroam/main.go](./freeroam/main.go)

This adds:

- spawn setup
- simple commands
- recurring tick logic

## Step 4: Add runtime context and error policies

- [context/main.go](./context/main.go)

Use this when you want to understand:

- lifecycle hooks
- event context decoration
- centralized runtime error handling

## Step 5: Study a fuller game mode

- [deathmatch/main.go](./deathmatch/main.go)

This example is useful once you want to see stateful gameplay logic with more handlers and more player-facing behavior.

## Step 6: Read the larger multi-file example

- [grandlarc/main.go](./grandlarc/main.go)
- [grandlarc/spawns.go](./grandlarc/spawns.go)

This is the most complex example in the repo. Save it for later.

## Internal test example

- [test/main.go](./test/main.go)

This exists mainly as a simple integration and smoke-test style example. It is useful for repo contributors, but it is not the recommended first example for new users.