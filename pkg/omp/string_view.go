package omp

import (
	"bytes"
	"encoding/json"
	"unsafe"
)

// BorrowedStringView is a borrowed, read-only view over callback-owned string data.
//
// The underlying memory is only guaranteed to remain valid for the duration of the
// callback dispatch. Use Clone or String when you need to retain the data.
type BorrowedStringView struct {
	data unsafe.Pointer
	len  int
}

// NewBorrowedStringView wraps borrowed callback-owned string memory.
func NewBorrowedStringView(data unsafe.Pointer, length int) BorrowedStringView {
	if data == nil || length <= 0 {
		return BorrowedStringView{}
	}
	return BorrowedStringView{data: data, len: length}
}

// Len returns the number of bytes in the borrowed view.
func (view BorrowedStringView) Len() int {
	return view.len
}

// IsEmpty reports whether the borrowed view is empty.
func (view BorrowedStringView) IsEmpty() bool {
	return view.data == nil || view.len == 0
}

// Clone copies the borrowed bytes into a Go string that is safe to keep.
func (view BorrowedStringView) Clone() string {
	if view.IsEmpty() {
		return ""
	}
	return string(view.bytes())
}

// String returns a copied Go string so the value is safe to log or retain.
func (view BorrowedStringView) String() string {
	return view.Clone()
}

// UnsafeString returns a zero-copy Go string view over the borrowed memory.
// The returned string must not outlive the callback dispatch that produced it.
func (view BorrowedStringView) UnsafeString() string {
	if view.IsEmpty() {
		return ""
	}
	return unsafe.String((*byte)(view.data), view.len)
}

// EqualString compares the borrowed bytes with a Go string without allocating.
func (view BorrowedStringView) EqualString(s string) bool {
	if view.len != len(s) {
		return false
	}
	if view.len == 0 {
		return true
	}
	return bytes.Equal(view.bytes(), unsafe.Slice(unsafe.StringData(s), len(s)))
}

// HasPrefix reports whether the borrowed bytes start with prefix without allocating.
func (view BorrowedStringView) HasPrefix(prefix string) bool {
	if len(prefix) > view.len {
		return false
	}
	if len(prefix) == 0 {
		return true
	}
	return bytes.Equal(view.bytes()[:len(prefix)], unsafe.Slice(unsafe.StringData(prefix), len(prefix)))
}

// MarshalJSON materializes the borrowed bytes as a copied JSON string.
func (view BorrowedStringView) MarshalJSON() ([]byte, error) {
	return json.Marshal(view.Clone())
}

func (view BorrowedStringView) bytes() []byte {
	if view.IsEmpty() {
		return nil
	}
	return unsafe.Slice((*byte)(view.data), view.len)
}
