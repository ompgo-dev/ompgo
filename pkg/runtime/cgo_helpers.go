package runtime

/*
#cgo CFLAGS: -I${SRCDIR}/include

#include <stdlib.h>
#include "ompcapi.h"
*/
import "C"

import "unsafe"

// CChar is an exported alias for C.char.
type CChar = C.char

// CFloat is an exported alias for C.float.
type CFloat = C.float

// CInt is an exported alias for C.int.
type CInt = C.int

// CUInt is an exported alias for C.uint.
type CUInt = C.uint

// CAPIStringView is an exported alias for the C API string view struct.
type CAPIStringView = C.struct_CAPIStringView

// CAPIStringBuffer is an exported alias for the C API string buffer struct.
type CAPIStringBuffer = C.struct_CAPIStringBuffer

// ComponentVersion is an exported alias for the C API component version struct.
type ComponentVersion = C.struct_ComponentVersion

// CString allocates a C string. Caller must free with FreeCString.
func CString(s string) *C.char {
	return C.CString(s)
}

// FreeCString frees a C string allocated by CString.
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}
