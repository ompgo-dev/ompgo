package runtime

import (
	"context"
	"fmt"
)

func invokeHandler[T any](ctx context.Context, stage string, event *T, fn func(context.Context, *T) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in %s: %v", stage, rec)
		}
	}()
	return fn(ctx, event)
}

func invokeBlockingHandler[T any](ctx context.Context, stage string, event *T, defaultAllowed bool, fn func(context.Context, *T) (bool, error)) (allowed bool, err error) {
	allowed = defaultAllowed
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in %s: %v", stage, rec)
		}
	}()
	return fn(ctx, event)
}

func invokeLifecycle(ctx context.Context, stage string, fn func(context.Context) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in %s: %v", stage, rec)
		}
	}()
	return fn(ctx)
}
