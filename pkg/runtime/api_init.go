package runtime

import (
	"context"
	"errors"
	"sync"
)

var (
	apiInitHooksMu sync.RWMutex
	apiInitHooks   []func(context.Context, *CAPI) error
)

// RegisterAPIInit registers a hook called after CAPI is Initialised.
func RegisterAPIInit(fn func(context.Context, *CAPI) error) {
	if fn == nil {
		return
	}
	apiInitHooksMu.Lock()
	apiInitHooks = append(apiInitHooks, fn)
	apiInitHooksMu.Unlock()
}

func callAPIInit(ctx context.Context, capi *CAPI) error {
	apiInitHooksMu.RLock()
	hooks := append([]func(context.Context, *CAPI) error(nil), apiInitHooks...)
	apiInitHooksMu.RUnlock()

	var errs []error
	for _, hook := range hooks {
		if err := invokeLifecycle(ctx, "Runtime.APIInit", func(ctx context.Context) error {
			return hook(ctx, capi)
		}); err != nil {
			if shouldBlockOnLifecycleError() {
				return err
			}
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
