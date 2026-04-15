package runtime

import (
	"unsafe"

	"github.com/ompgo-dev/ompgo/pkg/handle"
)

func handleFromPtr(p unsafe.Pointer) handle.Handle {
	return handle.Handle(uintptr(p))
}

func ptrFromHandle(h handle.Handle) unsafe.Pointer {
	return unsafe.Pointer(uintptr(h))
}
