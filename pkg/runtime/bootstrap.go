package runtime

import (
	"context"
	"sync"
)

// ComponentHandlers contains lifecycle callbacks invoked by the component.
type ComponentHandlers struct {
	OnReady func(context.Context) error
	OnReset func(context.Context) error
	OnFree  func(context.Context) error
}

// Config configures the runtime bootstrap process.
type Config struct {
	ComponentName    string
	ComponentVersion Version
	NewGamemode      func() Gamemode
	Handlers         ComponentHandlers
	Setup            func(context.Context) error
	BaseContext      context.Context
	ContextProvider  func() context.Context
	ContextDecorator EventContextDecorator
	EventError       EventErrorHandler
	LifecycleError   LifecycleErrorHandler
	EventErrorPolicy ErrorPolicy
	LifeErrorPolicy  ErrorPolicy
}

// Option is a decorator that configures Bootstrap.
type Option func(*Config)

var (
	handlersMu     sync.RWMutex
	activeHandlers = defaultHandlers()
)

// Bootstrap initializes the runtime configuration, registers the component,
// applies options, and runs Setup if provided.
func Bootstrap(opts ...Option) {
	cfg := Config{
		Handlers:         defaultHandlers(),
		BaseContext:      context.Background(),
		ContextProvider:  nil,
		ContextDecorator: nil,
		EventError:       nil,
		LifecycleError:   nil,
		EventErrorPolicy: ErrorPolicyContinue,
		LifeErrorPolicy:  ErrorPolicyContinue,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	gm := Gamemode(nil)
	if cfg.NewGamemode != nil {
		gm = cfg.NewGamemode()
	}

	if cfg.ComponentName != "" || (cfg.ComponentVersion != Version{}) {
		RegisterComponentWithVersion(cfg.ComponentName, cfg.ComponentVersion, gm)
	} else {
		RegisterComponent(registeredComponentName, gm)
	}

	setHandlers(cfg.Handlers)
	setEventDispatchConfig(eventDispatchConfig{
		baseContext:           cfg.BaseContext,
		contextProvider:       cfg.ContextProvider,
		decorator:             cfg.ContextDecorator,
		errorHandler:          cfg.EventError,
		lifecycleErrorHandler: cfg.LifecycleError,
		eventErrorPolicy:      cfg.EventErrorPolicy,
		lifecycleErrorPolicy:  cfg.LifeErrorPolicy,
	})

	if cfg.Setup != nil {
		ctx, cancel := newLifecycleContext("Bootstrap.Setup")
		defer cancel()
		if err := invokeLifecycle(ctx, "Bootstrap.Setup", cfg.Setup); err != nil {
			reportLifecycleError(ctx, "Bootstrap.Setup", err)
			if shouldBlockOnLifecycleError() {
				panic(err)
			}
		}
	}
}

// WithEventErrorPolicy sets how runtime event dispatch handles callback errors.
func WithEventErrorPolicy(policy ErrorPolicy) Option {
	return func(cfg *Config) {
		cfg.EventErrorPolicy = policy
	}
}

// WithLifecycleErrorPolicy sets how runtime lifecycle processing handles callback errors.
func WithLifecycleErrorPolicy(policy ErrorPolicy) Option {
	return func(cfg *Config) {
		cfg.LifeErrorPolicy = policy
	}
}

// WithBaseContext sets the base context used for event dispatch contexts.
func WithBaseContext(ctx context.Context) Option {
	return func(cfg *Config) {
		if ctx != nil {
			cfg.BaseContext = ctx
		}
	}
}

// WithContextProvider sets a provider used to supply the base context for each dispatched event.
func WithContextProvider(fn func() context.Context) Option {
	return func(cfg *Config) {
		cfg.ContextProvider = fn
	}
}

// WithEventContextDecorator sets a hook to decorate each per-event context before dispatch.
func WithEventContextDecorator(fn EventContextDecorator) Option {
	return func(cfg *Config) {
		cfg.ContextDecorator = fn
	}
}

// WithEventErrorHandler sets a hook invoked whenever an event handler returns an error.
func WithEventErrorHandler(fn EventErrorHandler) Option {
	return func(cfg *Config) {
		cfg.EventError = fn
	}
}

// WithLifecycleErrorHandler sets a hook invoked whenever a lifecycle callback returns an error.
func WithLifecycleErrorHandler(fn LifecycleErrorHandler) Option {
	return func(cfg *Config) {
		cfg.LifecycleError = fn
	}
}

// WithComponentName sets the component name.
func WithComponentName(name string) Option {
	return func(cfg *Config) {
		cfg.ComponentName = name
	}
}

// WithComponentVersion sets the component version.
func WithComponentVersion(version Version) Option {
	return func(cfg *Config) {
		cfg.ComponentVersion = version
	}
}

// WithGamemode sets the gamemode factory.
func WithGamemode(newGamemode func() Gamemode) Option {
	return func(cfg *Config) {
		cfg.NewGamemode = newGamemode
	}
}

// WithComponentHandlers sets the component lifecycle handlers.
func WithComponentHandlers(handlers ComponentHandlers) Option {
	return func(cfg *Config) {
		if handlers.OnReady != nil {
			cfg.Handlers.OnReady = handlers.OnReady
		}
		if handlers.OnReset != nil {
			cfg.Handlers.OnReset = handlers.OnReset
		}
		if handlers.OnFree != nil {
			cfg.Handlers.OnFree = handlers.OnFree
		}
	}
}

// WithOnReady sets the OnReady handler.
func WithOnReady(fn func(context.Context) error) Option {
	return func(cfg *Config) {
		if fn == nil {
			return
		}
		cfg.Handlers.OnReady = fn
	}
}

// WithOnReset sets the OnReset handler.
func WithOnReset(fn func(context.Context) error) Option {
	return func(cfg *Config) {
		if fn == nil {
			return
		}
		cfg.Handlers.OnReset = fn
	}
}

// WithOnFree sets the OnFree handler.
func WithOnFree(fn func(context.Context) error) Option {
	return func(cfg *Config) {
		if fn == nil {
			return
		}
		cfg.Handlers.OnFree = fn
	}
}

// WithSetup appends a setup function.
func WithSetup(fn func(context.Context) error) Option {
	return func(cfg *Config) {
		if fn == nil {
			return
		}
		prev := cfg.Setup
		cfg.Setup = func(ctx context.Context) error {
			if prev != nil {
				if err := prev(ctx); err != nil {
					return err
				}
			}
			return fn(ctx)
		}
	}
}

func defaultHandlers() ComponentHandlers {
	return ComponentHandlers{
		OnReady: func(context.Context) error { return Instance().Ready() },
		OnReset: func(context.Context) error { return nil },
		OnFree:  func(context.Context) error { return Instance().Unload() },
	}
}

func setHandlers(h ComponentHandlers) {
	handlersMu.Lock()
	activeHandlers = h
	handlersMu.Unlock()
}

func currentHandlers() ComponentHandlers {
	handlersMu.RLock()
	h := activeHandlers
	handlersMu.RUnlock()
	return h
}
