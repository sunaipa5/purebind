
//go:build darwin

package advanced

import "github.com/ebitengine/purego"

func init() {
	lib, err := purego.Dlopen("advanced.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	register(lib)
}
