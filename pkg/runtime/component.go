package runtime

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -ldl

#include <stdint.h>
#include <stdbool.h>

#include "ompcapi.h"

enum {
	OMPGO_LIFECYCLE_READY = 1,
	OMPGO_LIFECYCLE_RESET = 2,
	OMPGO_LIFECYCLE_FREE  = 3,
};

extern void OMPGODispatchLifecycle(int kind);

static struct OMPAPI_t g_ompgo_api;

static void ompgo_on_ready(void) {
	OMPGODispatchLifecycle(OMPGO_LIFECYCLE_READY);
}

static void ompgo_on_reset(void) {
	OMPGODispatchLifecycle(OMPGO_LIFECYCLE_RESET);
}

static void ompgo_on_free(void) {
	OMPGODispatchLifecycle(OMPGO_LIFECYCLE_FREE);
}

static void* ompgo_component_create(const char* name, struct ComponentVersion version) {
	if (!omp_initialize_capi(&g_ompgo_api) || !g_ompgo_api.Component.Create) {
		return NULL;
	}

	const uint64_t uid = 0x4f4d50474f000001ULL;
	return g_ompgo_api.Component.Create(
		uid,
		name,
		version,
		(void*)ompgo_on_ready,
		(void*)ompgo_on_reset,
		(void*)ompgo_on_free
	);
}
*/
import "C"

import (
	"log"
	"unsafe"
)

//export OMPGODispatchLifecycle
func OMPGODispatchLifecycle(kind C.int) {
	switch int(kind) {
	case 1:
		OnReady()
	case 2:
		OnReset()
	case 3:
		OnFree()
	}
}

// Version describes a component version for registration.
type Version struct {
	Major  uint8
	Minor  uint8
	Patch  uint8
	Prerel uint16
}

var (
	registeredComponentName    = "ompgo"
	registeredComponentVersion = Version{Major: 0, Minor: 1, Patch: 0, Prerel: 0}
	registeredComponentSet     = false
)

// RegisterComponent registers the component name and gamemode.
// This should be called from init() in the user component.
func RegisterComponent(name string, gm Gamemode) {
	RegisterComponentWithVersion(name, registeredComponentVersion, gm)
}

// RegisterComponentWithVersion registers the component name, version, and gamemode.
func RegisterComponentWithVersion(name string, version Version, gm Gamemode) {
	if name != "" {
		registeredComponentName = name
	}
	registeredComponentVersion = version
	registeredComponentSet = true
	if gm != nil {
		Instance().SetGamemode(gm)
	}
}

// Init initializes the runtime and loads the C API.
func Init(capi unsafe.Pointer) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ompgo] PANIC in Init: %v", r)
		}
	}()

	rt := Instance()
	if err := rt.Load(capi); err != nil {
		log.Printf("[ompgo] Error during Load: %v", err)
	}
}

// ComponentEntryPoint creates the open.mp component instance.
//
//export ComponentEntryPoint
func ComponentEntryPoint() unsafe.Pointer {
	if !registeredComponentSet {
		registeredComponentSet = true
	}
	Init(nil)

	cname := CString(registeredComponentName)
	defer FreeCString(cname)

	version := C.struct_ComponentVersion{
		major:  C.uint8_t(registeredComponentVersion.Major),
		minor:  C.uint8_t(registeredComponentVersion.Minor),
		patch:  C.uint8_t(registeredComponentVersion.Patch),
		prerel: C.uint16_t(registeredComponentVersion.Prerel),
	}

	return unsafe.Pointer(C.ompgo_component_create(cname, version))
}

// ComponentCleanup releases runtime resources.
//
//export ComponentCleanup
func ComponentCleanup() {
	OnFree()
}

// OnFree handles component free callback.
func OnFree() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ompgo] PANIC in OnFree: %v", r)
		}
	}()

	ctx, cancel := newLifecycleContext("Component.OnFree")
	defer cancel()
	if h := currentHandlers().OnFree; h != nil {
		if err := invokeLifecycle(ctx, "Component.OnFree", h); err != nil {
			reportLifecycleError(ctx, "Component.OnFree", err)
		}
	}
}

// OnReady handles component ready callback.
func OnReady() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ompgo] PANIC in OnReady: %v", r)
		}
	}()

	ctx, cancel := newLifecycleContext("Component.OnReady")
	defer cancel()
	if h := currentHandlers().OnReady; h != nil {
		if err := invokeLifecycle(ctx, "Component.OnReady", h); err != nil {
			reportLifecycleError(ctx, "Component.OnReady", err)
		}
	}
}

// OnReset handles component reset callback.
func OnReset() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ompgo] PANIC in OnReset: %v", r)
		}
	}()

	ctx, cancel := newLifecycleContext("Component.OnReset")
	defer cancel()
	if h := currentHandlers().OnReset; h != nil {
		if err := invokeLifecycle(ctx, "Component.OnReset", h); err != nil {
			reportLifecycleError(ctx, "Component.OnReset", err)
		}
	}
}
