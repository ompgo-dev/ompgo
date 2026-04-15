package runtime

// CAPIStringViewToString converts a CAPI string view to a Go string.
func CAPIStringViewToString(view *CAPIStringView) string {
	if view == nil {
		return ""
	}
	return stringFromView(*view)
}

// CAPIStringBufferToString converts a CAPI string buffer to a Go string.
func CAPIStringBufferToString(buf *CAPIStringBuffer) string {
	if buf == nil {
		return ""
	}
	return stringFromBuffer(*buf)
}
