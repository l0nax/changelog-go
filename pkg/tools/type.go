package tools

import "unsafe"

// ByteSlice2String converts a byte slice to a string in a performant way.
func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
