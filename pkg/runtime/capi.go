package runtime

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -ldl

#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

#include "ompcapi.h"

bool capi_register_all_handlers(void);

struct OMPAPI_t g_ompapi;
bool g_capi_initialised = false;

static bool capi_init() {
    if (g_capi_initialised) {
        return true;
    }
    g_capi_initialised = omp_initialize_capi(&g_ompapi);
    return g_capi_initialised;
}

static bool capi_is_initialised() {
    return g_capi_initialised;
}

static void capi_cleanup() {
    g_capi_initialised = false;
}

static bool capi_register_basic_handlers() {
    return capi_register_all_handlers();
}
*/
import "C"

import "fmt"

// CAPI provides access to the open.mp C API loaded from $CAPI.so
type CAPI struct {
	Initialised bool
}

// Initialise loads the C API and registers default event handlers.
func (c *CAPI) Initialise() error {
	if c == nil {
		return fmt.Errorf("capi: nil receiver")
	}
	if c.Initialised {
		return nil
	}
	if C.capi_init() == C.bool(false) {
		return fmt.Errorf("capi: failed to initialise open.mp C API ($CAPI.so)")
	}
	if C.capi_register_basic_handlers() == C.bool(false) {
		return fmt.Errorf("capi: failed to register event handlers")
	}
	c.Initialised = true
	return nil
}

// IsInitialised returns whether the C API has been Initialised.
func (c *CAPI) IsInitialised() bool {
	return c != nil && c.Initialised && C.capi_is_initialised() == C.bool(true)
}

// Cleanup cleans up the C API resources.
func (c *CAPI) Cleanup() {
	if c == nil {
		return
	}
	c.Initialised = false
	C.capi_cleanup()
}
