
package advanced

import (
	"unsafe"
	"github.com/ebitengine/purego"
)

var handle uintptr

var (
	complexFunction func(int32,float32,unsafe.Pointer,unsafe.Pointer) uintptr
)

func register(libHandle uintptr) {
	handle = libHandle
	purego.RegisterLibFunc(&complexFunction, handle, "ComplexFunction")
}
