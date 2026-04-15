package runtime

import (
	"sync"
	"sync/atomic"
)

type handlerRegistration[F any] struct {
	id int64
	fn F
}

type handlerSnapshot[F any] []handlerRegistration[F]

func loadHandlerSnapshot[F any](ptr *atomic.Pointer[handlerSnapshot[F]]) handlerSnapshot[F] {
	current := ptr.Load()
	if current == nil {
		return nil
	}
	return *current
}

func registerHandler[F any](mu *sync.Mutex, seq *atomic.Int64, ptr *atomic.Pointer[handlerSnapshot[F]], handler F) func() {
	id := seq.Add(1)

	mu.Lock()
	current := loadHandlerSnapshot(ptr)
	next := make(handlerSnapshot[F], 0, len(current)+1)
	next = append(next, current...)
	next = append(next, handlerRegistration[F]{id: id, fn: handler})
	ptr.Store(&next)
	mu.Unlock()

	return func() {
		unregisterHandler(mu, ptr, id)
	}
}

func unregisterHandler[F any](mu *sync.Mutex, ptr *atomic.Pointer[handlerSnapshot[F]], id int64) {
	mu.Lock()
	defer mu.Unlock()

	current := loadHandlerSnapshot(ptr)
	if len(current) == 0 {
		return
	}

	next := make(handlerSnapshot[F], 0, len(current))
	for _, registered := range current {
		if registered.id == id {
			continue
		}
		next = append(next, registered)
	}

	if len(next) == 0 {
		ptr.Store(nil)
		return
	}
	ptr.Store(&next)
}
