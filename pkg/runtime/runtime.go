// Package runtime provides the core Go runtime for the open.mp component.
// This package manages the component lifecycle, C API integration, and
// event dispatching.
package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ompgo-dev/ompgo/pkg/omp"
)

// Gamemode is the interface that all gamemodes must implement.
// This is defined here to avoid circular imports.
// It embeds omp.EventHandler to receive all events and adds OnLoad lifecycle hook
type Gamemode interface {
	omp.EventHandler
	OnLoad(ctx context.Context) error
}

// State represents the current runtime state
type State int32

const (
	StateUninitialised State = iota
	StateLoading
	StateLoaded
	StateReady
	StateUnloading
)

// Runtime manages the Go component lifecycle
type Runtime struct {
	state          atomic.Int32
	mu             sync.Mutex
	activeGamemode atomic.Pointer[gamemodeSnapshot]
	// C API reference
	capi *CAPI
}

type gamemodeSnapshot struct {
	gamemode Gamemode
	loaded   bool
}

func (r *Runtime) loadGamemodeSnapshot() gamemodeSnapshot {
	if snapshot := r.activeGamemode.Load(); snapshot != nil {
		return *snapshot
	}
	return gamemodeSnapshot{}
}

func (r *Runtime) storeGamemodeSnapshot(gm Gamemode, loaded bool) {
	if gm == nil {
		r.activeGamemode.Store(nil)
		return
	}

	snapshot := &gamemodeSnapshot{
		gamemode: gm,
		loaded:   loaded,
	}
	r.activeGamemode.Store(snapshot)
}

func (r *Runtime) currentGamemode() Gamemode {
	snapshot := r.loadGamemodeSnapshot()

	if r.State() == StateReady && !snapshot.loaded {
		return nil
	}
	return snapshot.gamemode
}

// Global runtime instance
var (
	instance     *Runtime
	instanceOnce sync.Once
)

// Instance returns the singleton runtime instance
func Instance() *Runtime {
	instanceOnce.Do(func() {
		instance = &Runtime{
			capi: &CAPI{},
		}
	})
	return instance
}

// State returns the current runtime state
func (r *Runtime) State() State {
	return State(r.state.Load())
}

// setState atomically sets the runtime state
func (r *Runtime) setState(state State) {
	r.state.Store(int32(state))
}

// SetGamemode registers the gamemode implementation
func (r *Runtime) SetGamemode(gm Gamemode) {
	shouldLoad := false

	r.mu.Lock()
	r.storeGamemodeSnapshot(gm, false)

	// Call OnLoad only after the runtime is ready
	if gm != nil && r.State() == StateReady {
		shouldLoad = true
	}
	r.mu.Unlock()

	if !shouldLoad {
		return
	}

	ctx, cancel := newLifecycleContext("Gamemode.OnLoad")
	defer cancel()
	if err := invokeLifecycle(ctx, "Gamemode.OnLoad", gm.OnLoad); err != nil {
		reportLifecycleError(ctx, "Gamemode.OnLoad", err)
		if shouldBlockOnLifecycleError() {
			panic(err)
		}
		return
	}

	r.mu.Lock()
	if current := r.loadGamemodeSnapshot(); current.gamemode == gm {
		r.storeGamemodeSnapshot(gm, true)
	}
	r.mu.Unlock()
}

// Gamemode returns the current gamemode implementation
func (r *Runtime) Gamemode() Gamemode {
	return r.loadGamemodeSnapshot().gamemode
}

// CAPI returns the C API interface
func (r *Runtime) CAPI() *CAPI {
	return r.capi
}

// Load is called when the component is loaded by the server.
// It receives the core API pointer and initializes the runtime.
//
// Parameters:
//   - core: Pointer to the open.mp core API interface
//
// Returns:
//   - error: If initialization fails
func (r *Runtime) Load(core unsafe.Pointer) error {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("[Runtime] PANIC in Load: %v\n", rec)
		}
	}()

	if r.State() != StateUninitialised {
		fmt.Printf("[Runtime] Already loaded (state: %d)\n", r.State())
		return fmt.Errorf("component already loaded")
	}

	r.setState(StateLoading)

	// Initialise CAPI via $CAPI.so
	if err := r.capi.Initialise(); err != nil {
		fmt.Printf("[Runtime] C API initialisation failed: %v\n", err)
		r.setState(StateUninitialised)
		return err
	}
	ctx, cancel := newLifecycleContext("Runtime.Load")
	defer cancel()
	if err := callAPIInit(ctx, r.capi); err != nil {
		reportLifecycleError(ctx, "Runtime.Load", err)
		r.setState(StateUninitialised)
		return err
	}

	r.setState(StateLoaded)
	return nil
}

// Ready is called when the server is ready
func (r *Runtime) Ready() error {
	if r.State() != StateLoaded {
		return nil
	}

	r.setState(StateReady)

	snapshot := r.loadGamemodeSnapshot()
	if snapshot.gamemode != nil && !snapshot.loaded {
		ctx, cancel := newLifecycleContext("Gamemode.OnLoad")
		defer cancel()
		if err := invokeLifecycle(ctx, "Gamemode.OnLoad", snapshot.gamemode.OnLoad); err != nil {
			reportLifecycleError(ctx, "Gamemode.OnLoad", err)
			return err
		}
		r.mu.Lock()
		if current := r.loadGamemodeSnapshot(); current.gamemode == snapshot.gamemode {
			r.storeGamemodeSnapshot(snapshot.gamemode, true)
		}
		r.mu.Unlock()
	}
	return nil
}

// Tick is called every server tick
func (r *Runtime) Tick(elapsed int) {
	if r.State() != StateReady {
		return
	}
}

// Unload is called when the component is being unloaded
func (r *Runtime) Unload() error {
	if r.State() == StateUninitialised {
		return nil
	}

	r.setState(StateUnloading)

	r.capi.Cleanup()

	r.setState(StateUninitialised)
	return nil
}
