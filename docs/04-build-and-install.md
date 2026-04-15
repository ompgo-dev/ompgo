# Step 4: Build and install it

This step turns your Go project into a component that open.mp can load.

## Build the shared library

From your project root:

```bash
CGO_ENABLED=1 go build -buildmode=c-shared -o ompgo_mygamemode.so .
```

Adjust the output filename to match your component name.

## Install it into your server

1. Copy the generated `.so` file into your server's `components` folder.
2. Make sure `$CAPI.so` is in the same folder.
3. Start the server.
4. Check the server logs to confirm your component loaded.

## What a successful first run looks like

You should see:

- the server loading your shared library
- your `OnLoad` log line
- player-facing behavior from your first event handler, if you added one

## Common setup mistakes

- `CGO_ENABLED` was off during build
- the built filename does not match the component you intended to load
- `$CAPI.so` is missing from the server's `components` folder
- the server is loading an older build from a different path

## Move on when

Move on once the server loads your component and your startup log appears.

Next: [Step 5: Add commands and modules](./05-commands-and-modules.md)