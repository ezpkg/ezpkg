package unsafez // import "ezpkg.io/unsafez"

import (
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
