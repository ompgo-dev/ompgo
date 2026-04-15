# Step 6: Add context and error handling

This step is for when the basic component works and you want better diagnostics and safer runtime behavior.

## Why this matters

As your component grows, you usually need to answer three questions:

1. Which event caused this log or error?
2. Should a handler error block the event or be ignored?
3. How do I attach request-like metadata to each event?

`pkg/runtime` gives you tools for all three.

## Add lifecycle hooks and policies

```go
func init() {
	runtime.Bootstrap(
		runtime.WithComponentName("ompgo_context_demo"),
		runtime.WithGamemode(NewGamemode),
		runtime.WithSetup(func(ctx context.Context) error {
			return nil
		}),
		runtime.WithOnReady(func(ctx context.Context) error {
			return nil
		}),
		runtime.WithOnFree(func(ctx context.Context) error {
			return nil
		}),
		runtime.WithEventErrorPolicy(runtime.ErrorPolicyBlockOnError),
		runtime.WithLifecycleErrorPolicy(runtime.ErrorPolicyContinue),
	)
}
```

Use this when you need explicit setup or teardown work and predictable failure behavior.

## Decorate the event context

You can add metadata to each event context before your handlers see it.

```go
runtime.WithEventContextDecorator(func(ctx context.Context, eventName string, event any) context.Context {
	_ = event
	return runtime.WithRequestID(ctx, eventName+"-request")
})
```

Inside handlers, you can read values back with:

- `runtime.RequestIDFromContext(ctx)`
- `runtime.EventNameFromContext(ctx)`
- `runtime.EventStartedAtFromContext(ctx)`

## Centralize error reporting

```go
runtime.WithEventErrorHandler(func(ctx context.Context, eventName string, event any, err error) {
	_ = ctx
	_ = event
	log.Printf("event error in %s: %v", eventName, err)
})
```

This keeps error handling out of your gameplay code and gives you a single reporting path.

## See a full example

The best reference for this step is [examples/context/main.go](../examples/context/main.go).

## Where to go next

After this step, you have the core pieces most components need.

From here, pick the example that matches what you want to learn next:

- [examples/basic/main.go](../examples/basic/main.go) for a compact baseline
- [examples/freeroam/main.go](../examples/freeroam/main.go) for a gameplay-oriented flow
- [examples/grandlarc/main.go](../examples/grandlarc/main.go) for a larger example

For the full example learning path, use [examples/README.md](../examples/README.md).