package runtime

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/ompgo-dev/ompgo/pkg/omp"
)

type testGamemode struct {
	omp.BaseEventHandler
	onLoad              func(context.Context) error
	onPlayerCommandText func(context.Context, *omp.PlayerCommandTextEvent) (bool, error)
}

func (g *testGamemode) OnLoad(ctx context.Context) error {
	if g.onLoad != nil {
		return g.onLoad(ctx)
	}
	return nil
}

func (g *testGamemode) OnPlayerCommandText(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
	if g.onPlayerCommandText != nil {
		return g.onPlayerCommandText(ctx, event)
	}
	return true, nil
}

func TestOnPlayerCommandTextBehavior(t *testing.T) {
	tests := []struct {
		name            string
		registerHandler func() (func(), bool)
		gmFn            func(context.Context, *omp.PlayerCommandTextEvent) (bool, error)
		policy          ErrorPolicy
		wantAllowed     bool
		wantErrors      int
		wantGMCalled    bool
	}{
		{
			name: "handler_false_blocks",
			registerHandler: func() (func(), bool) {
				return RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
					return false, nil
				}), false
			},
			gmFn:         func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil },
			policy:       ErrorPolicyContinue,
			wantAllowed:  false,
			wantErrors:   0,
			wantGMCalled: false,
		},
		{
			name: "handler_panic_reports_and_continues",
			registerHandler: func() (func(), bool) {
				return RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
					panic("handler boom")
				}), true
			},
			gmFn:         func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil },
			policy:       ErrorPolicyContinue,
			wantAllowed:  true,
			wantErrors:   1,
			wantGMCalled: true,
		},
		{
			name:            "gm_panic_reports_default_return",
			registerHandler: func() (func(), bool) { return func() {}, false },
			gmFn: func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
				panic("gm boom")
			},
			policy:       ErrorPolicyContinue,
			wantAllowed:  true,
			wantErrors:   1,
			wantGMCalled: true,
		},
		{
			name: "handler_error_blocks_when_policy_block",
			registerHandler: func() (func(), bool) {
				return RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
					return true, errors.New("handler error")
				}), true
			},
			gmFn:         func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil },
			policy:       ErrorPolicyBlockOnError,
			wantAllowed:  false,
			wantErrors:   1,
			wantGMCalled: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			oldCfg := *currentEventDispatchConfig()
			t.Cleanup(func() { setEventDispatchConfig(oldCfg) })

			errorCount := 0
			cfg := oldCfg
			cfg.errorHandler = func(context.Context, string, any, error) {
				errorCount++
			}
			cfg.eventErrorPolicy = tc.policy
			setEventDispatchConfig(cfg)

			gmCalled := false
			rt := &Runtime{capi: &CAPI{}}
			rt.storeGamemodeSnapshot(&testGamemode{
				onPlayerCommandText: func(ctx context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
					gmCalled = true
					return tc.gmFn(ctx, event)
				},
			}, true)

			unregister, _ := tc.registerHandler()
			t.Cleanup(unregister)

			allowed := rt.OnPlayerCommandText(nil, omp.NewBorrowedStringView(unsafe.Pointer(unsafe.StringData("/x")), len("/x")))
			if allowed != tc.wantAllowed {
				t.Fatalf("allowed=%v want=%v", allowed, tc.wantAllowed)
			}
			if gmCalled != tc.wantGMCalled {
				t.Fatalf("gmCalled=%v want=%v", gmCalled, tc.wantGMCalled)
			}
			if errorCount != tc.wantErrors {
				t.Fatalf("errorCount=%d want=%d", errorCount, tc.wantErrors)
			}
		})
	}
}

func TestOnPlayerCommandTextUsesBorrowedStringView(t *testing.T) {
	rt := &Runtime{capi: &CAPI{}}
	gmCalled := false
	rt.storeGamemodeSnapshot(&testGamemode{
		onPlayerCommandText: func(_ context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
			gmCalled = true
			if !event.Command.EqualString("/vehicle repair 411") {
				t.Fatalf("command mismatch: got %q", event.Command.Clone())
			}
			if !event.Command.HasPrefix("/vehicle") {
				t.Fatal("command missing expected prefix")
			}
			if event.Command.Len() != len("/vehicle repair 411") {
				t.Fatalf("command len=%d", event.Command.Len())
			}
			return true, nil
		},
	}, true)

	unregister := RegisterOnPlayerCommandText(func(_ context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
		if !event.Command.EqualString("/vehicle repair 411") {
			t.Fatalf("handler command mismatch: got %q", event.Command.Clone())
		}
		return true, nil
	})
	defer unregister()

	viewPtr := CString("/vehicle repair 411")
	defer FreeCString(viewPtr)

	allowed := rt.OnPlayerCommandText(nil, borrowedStringViewFromCAPIStringView(CAPIStringView{
		len:  CUInt(len("/vehicle repair 411")),
		data: viewPtr,
	}))
	if !allowed {
		t.Fatal("allowed = false, want true")
	}
	if !gmCalled {
		t.Fatal("gamemode was not called")
	}
}

func TestReadyOnLoadErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		onLoad     func(context.Context) error
		wantErr    bool
		wantLoaded bool
		wantReport int
	}{
		{
			name:       "onload_success",
			onLoad:     func(context.Context) error { return nil },
			wantErr:    false,
			wantLoaded: true,
			wantReport: 0,
		},
		{
			name:       "onload_error_reported",
			onLoad:     func(context.Context) error { return errors.New("onload failed") },
			wantErr:    true,
			wantLoaded: false,
			wantReport: 1,
		},
		{
			name: "onload_panic_reported",
			onLoad: func(context.Context) error {
				panic("onload panic")
			},
			wantErr:    true,
			wantLoaded: false,
			wantReport: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			oldCfg := *currentEventDispatchConfig()
			t.Cleanup(func() { setEventDispatchConfig(oldCfg) })

			lifecycleReports := 0
			cfg := oldCfg
			cfg.lifecycleErrorHandler = func(context.Context, string, error) {
				lifecycleReports++
			}
			setEventDispatchConfig(cfg)

			rt := &Runtime{capi: &CAPI{}}
			rt.storeGamemodeSnapshot(&testGamemode{onLoad: tc.onLoad}, false)
			rt.setState(StateLoaded)

			err := rt.Ready()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if snapshot := rt.loadGamemodeSnapshot(); snapshot.loaded != tc.wantLoaded {
				t.Fatalf("loaded=%v want=%v", snapshot.loaded, tc.wantLoaded)
			}
			if lifecycleReports != tc.wantReport {
				t.Fatalf("lifecycleReports=%d want=%d", lifecycleReports, tc.wantReport)
			}
		})
	}
}

func TestSetGamemodeDoesNotExposeReadyGamemodeBeforeOnLoadCompletes(t *testing.T) {
	t.Parallel()

	rt := &Runtime{capi: &CAPI{}}
	rt.setState(StateReady)

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	gm := &testGamemode{
		onLoad: func(context.Context) error {
			close(started)
			<-release
			return nil
		},
	}

	go func() {
		rt.SetGamemode(gm)
		close(done)
	}()

	<-started

	readerDone := make(chan struct{})
	go func() {
		_ = rt.currentGamemode()
		close(readerDone)
	}()

	select {
	case <-readerDone:
	case <-time.After(time.Second):
		t.Fatal("currentGamemode blocked while OnLoad was running")
	}

	if got := rt.currentGamemode(); got != nil {
		t.Fatalf("currentGamemode() = %#v, want nil before OnLoad completes", got)
	}

	close(release)
	<-done

	if got := rt.currentGamemode(); got != gm {
		t.Fatalf("currentGamemode() = %#v, want %#v after OnLoad completes", got, gm)
	}
	if !rt.loadGamemodeSnapshot().loaded {
		t.Fatal("runtime.loaded = false, want true after OnLoad completes")
	}
}

func TestSetGamemodeReplacesLoadedGamemodeState(t *testing.T) {
	t.Parallel()

	rt := &Runtime{capi: &CAPI{}}
	rt.setState(StateReady)
	rt.storeGamemodeSnapshot(&testGamemode{}, true)

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	next := &testGamemode{
		onLoad: func(context.Context) error {
			close(started)
			<-release
			return nil
		},
	}

	go func() {
		rt.SetGamemode(next)
		close(done)
	}()

	<-started

	if got := rt.currentGamemode(); got != nil {
		t.Fatalf("currentGamemode() = %#v, want nil while replacement OnLoad is running", got)
	}

	close(release)
	<-done

	if got := rt.currentGamemode(); got != next {
		t.Fatalf("currentGamemode() = %#v, want %#v after replacement OnLoad completes", got, next)
	}
}

func TestCallAPIInitAggregatesErrors(t *testing.T) {
	apiInitHooksMu.Lock()
	oldHooks := append([]func(context.Context, *CAPI) error(nil), apiInitHooks...)
	apiInitHooks = nil
	apiInitHooksMu.Unlock()
	t.Cleanup(func() {
		apiInitHooksMu.Lock()
		apiInitHooks = oldHooks
		apiInitHooksMu.Unlock()
	})

	RegisterAPIInit(func(context.Context, *CAPI) error { return errors.New("first failure") })
	RegisterAPIInit(func(context.Context, *CAPI) error { panic("panic failure") })
	RegisterAPIInit(func(context.Context, *CAPI) error { return errors.New("third failure") })

	err := callAPIInit(context.Background(), &CAPI{})
	if err == nil {
		t.Fatalf("expected aggregated error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "first failure") || !strings.Contains(msg, "third failure") || !strings.Contains(msg, "panic in Runtime.APIInit") {
		t.Fatalf("unexpected aggregated error: %v", err)
	}
}

func TestCallAPIInitStopsOnFirstErrorWhenPolicyBlock(t *testing.T) {
	apiInitHooksMu.Lock()
	oldHooks := append([]func(context.Context, *CAPI) error(nil), apiInitHooks...)
	apiInitHooks = nil
	apiInitHooksMu.Unlock()
	t.Cleanup(func() {
		apiInitHooksMu.Lock()
		apiInitHooks = oldHooks
		apiInitHooksMu.Unlock()
	})

	oldCfg := *currentEventDispatchConfig()
	t.Cleanup(func() { setEventDispatchConfig(oldCfg) })
	cfg := oldCfg
	cfg.lifecycleErrorPolicy = ErrorPolicyBlockOnError
	setEventDispatchConfig(cfg)

	calls := 0
	RegisterAPIInit(func(context.Context, *CAPI) error {
		calls++
		return errors.New("first failure")
	})
	RegisterAPIInit(func(context.Context, *CAPI) error {
		calls++
		return errors.New("second failure")
	})

	err := callAPIInit(context.Background(), &CAPI{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if calls != 1 {
		t.Fatalf("calls=%d want=1", calls)
	}
}
