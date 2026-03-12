package advanced

import "unsafe"

func ComplexFunction(p0 int32, values float32, name unsafe.Pointer, tex unsafe.Pointer) uintptr {
	return complexFunction(p0, values, name, tex)
}
