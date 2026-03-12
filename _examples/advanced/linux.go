
//go:build linux

package advanced

import "github.com/ebitengine/purego"

func init() {
	lib, err := purego.Dlopen("advanced.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	register(lib)
}
