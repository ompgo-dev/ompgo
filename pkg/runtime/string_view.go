package runtime

import (
	"unsafe"

	"github.com/ompgo-dev/ompgo/pkg/omp"
)

func borrowedStringViewFromCAPIStringView(view CAPIStringView) omp.BorrowedStringView {
	return omp.NewBorrowedStringView(unsafe.Pointer(view.data), int(view.len))
}
