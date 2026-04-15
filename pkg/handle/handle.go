package handle

// Handle is an opaque reference to an open.mp entity (C pointer).
//
// It's intentionally not an unsafe.Pointer in the public API so end-user code
// doesn't need to import "unsafe" to work with handles.
type Handle uintptr

func (h Handle) Valid() bool { return h != 0 }
