package runtime

import (
	"context"
	"testing"
	"unsafe"

	"github.com/ompgo-dev/ompgo/pkg/omp"
)

func benchmarkRuntime() *Runtime {
	rt := &Runtime{capi: &CAPI{}}
	rt.setState(StateReady)
	rt.storeGamemodeSnapshot(&testGamemode{}, true)
	return rt
}

func benchmarkGlobalRuntime(b *testing.B, gm Gamemode) *Runtime {
	b.Helper()

	rt := Instance()
	prevState := rt.State()
	prevSnapshot := rt.loadGamemodeSnapshot()
	b.Cleanup(func() {
		rt.setState(prevState)
		rt.storeGamemodeSnapshot(prevSnapshot.gamemode, prevSnapshot.loaded)
	})

	rt.setState(StateReady)
	rt.storeGamemodeSnapshot(gm, true)
	return rt
}

func composeUnregisters(unregisters ...func()) func() {
	return func() {
		for i := len(unregisters) - 1; i >= 0; i-- {
			if unregisters[i] != nil {
				unregisters[i]()
			}
		}
	}
}

func benchmarkStringView(b *testing.B, text string) CAPIStringView {
	b.Helper()
	ptr := CString(text)
	b.Cleanup(func() { FreeCString(ptr) })
	return CAPIStringView{
		len:  CUInt(len(text)),
		data: ptr,
	}
}

func BenchmarkRuntimeOnTick(b *testing.B) {
	benchmarks := []struct {
		name     string
		register func() func()
	}{
		{
			name: "gamemode_only",
			register: func() func() {
				return func() {}
			},
		},
		{
			name: "with_handler",
			register: func() func() {
				return RegisterOnTick(func(context.Context, *omp.TickEvent) error { return nil })
			},
		},
		{
			name: "with_4_handlers",
			register: func() func() {
				return composeUnregisters(
					RegisterOnTick(func(context.Context, *omp.TickEvent) error { return nil }),
					RegisterOnTick(func(context.Context, *omp.TickEvent) error { return nil }),
					RegisterOnTick(func(context.Context, *omp.TickEvent) error { return nil }),
					RegisterOnTick(func(context.Context, *omp.TickEvent) error { return nil }),
				)
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			rt := benchmarkRuntime()

			unregister := bm.register()
			b.Cleanup(unregister)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rt.OnTick(16)
			}
		})
	}
}

func BenchmarkRuntimeOnPlayerUpdate(b *testing.B) {
	benchmarks := []struct {
		name     string
		register func() func()
	}{
		{
			name: "gamemode_only",
			register: func() func() {
				return func() {}
			},
		},
		{
			name: "with_handler",
			register: func() func() {
				return RegisterOnPlayerUpdate(func(context.Context, *omp.PlayerUpdateEvent) (bool, error) {
					return true, nil
				})
			},
		},
		{
			name: "with_4_handlers",
			register: func() func() {
				return composeUnregisters(
					RegisterOnPlayerUpdate(func(context.Context, *omp.PlayerUpdateEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerUpdate(func(context.Context, *omp.PlayerUpdateEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerUpdate(func(context.Context, *omp.PlayerUpdateEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerUpdate(func(context.Context, *omp.PlayerUpdateEvent) (bool, error) { return true, nil }),
				)
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			rt := benchmarkRuntime()

			unregister := bm.register()
			b.Cleanup(unregister)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if !rt.OnPlayerUpdate(nil) {
					b.Fatal("OnPlayerUpdate unexpectedly blocked")
				}
			}
		})
	}
}

func BenchmarkRuntimeOnPlayerCommandText(b *testing.B) {
	benchmarks := []struct {
		name     string
		register func() func()
	}{
		{
			name: "gamemode_only",
			register: func() func() {
				return func() {}
			},
		},
		{
			name: "with_handler",
			register: func() func() {
				return RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
					return true, nil
				})
			},
		},
		{
			name: "with_4_handlers",
			register: func() func() {
				return composeUnregisters(
					RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil }),
					RegisterOnPlayerCommandText(func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) { return true, nil }),
				)
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			rt := benchmarkRuntime()
			rt.storeGamemodeSnapshot(&testGamemode{
				onPlayerCommandText: func(context.Context, *omp.PlayerCommandTextEvent) (bool, error) {
					return true, nil
				},
			}, true)

			unregister := bm.register()
			b.Cleanup(unregister)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if !rt.OnPlayerCommandText(nil, omp.NewBorrowedStringView(unsafe.Pointer(unsafe.StringData("/vehicle repair 411")), len("/vehicle repair 411"))) {
					b.Fatal("OnPlayerCommandText unexpectedly blocked")
				}
			}
		})
	}
}

func BenchmarkStringFromView(b *testing.B) {
	view := benchmarkStringView(b, "/vehicle repair 411")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stringFromView(view)
	}
}

func BenchmarkOMPGOOnPlayerCommandText(b *testing.B) {
	view := benchmarkStringView(b, "/vehicle repair 411")
	benchmarkGlobalRuntime(b, &testGamemode{
		onPlayerCommandText: func(_ context.Context, event *omp.PlayerCommandTextEvent) (bool, error) {
			if !event.Command.HasPrefix("/vehicle") {
				b.Fatal("gamemode received unexpected command")
			}
			return true, nil
		},
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = OMPGO_OnPlayerCommandText(nil, view)
	}
}

func BenchmarkOMPGOOnPlayerText(b *testing.B) {
	view := benchmarkStringView(b, "hello world")
	benchmarkGlobalRuntime(b, &testGamemode{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = OMPGO_OnPlayerText(nil, view)
	}
}

func BenchmarkOMPGOOnIncomingConnection(b *testing.B) {
	view := benchmarkStringView(b, "203.0.113.42")
	benchmarkGlobalRuntime(b, &testGamemode{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OMPGO_OnIncomingConnection(nil, view, CInt(7777))
	}
}

func BenchmarkOMPGOOnConsoleText(b *testing.B) {
	command := benchmarkStringView(b, "vehicle")
	parameters := benchmarkStringView(b, "repair 411")
	benchmarkGlobalRuntime(b, &testGamemode{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = OMPGO_OnConsoleText(command, parameters)
	}
}
