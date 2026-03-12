
//go:build windows

package advanced

import "syscall"

func init() {
	libvlc, err := syscall.LoadLibrary("advanced.dll")
	if err != nil {
		panic(err)
	}

	register(uintptr(libvlc))
}
