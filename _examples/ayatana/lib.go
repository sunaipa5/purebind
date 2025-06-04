
package ayatana

import (
	"unsafe"
	"github.com/ebitengine/purego"
)

var handle uintptr

var (
	app_indicator_set_status func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_attention_icon_full func(unsafe.Pointer,unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_menu func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_icon_full func(unsafe.Pointer,unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_label func(unsafe.Pointer,unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_icon_theme_path func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_ordering_index func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_secondary_activate_target func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_set_title func(unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
	app_indicator_get_category func(unsafe.Pointer) unsafe.Pointer
	app_indicator_get_status func(unsafe.Pointer) unsafe.Pointer
	app_indicator_get_ordering_index func(unsafe.Pointer) unsafe.Pointer
	app_indicator_build_menu_from_desktop func(unsafe.Pointer,unsafe.Pointer,unsafe.Pointer) unsafe.Pointer
)

func Register(libHandle uintptr) {
	handle = libHandle
	purego.RegisterLibFunc(&app_indicator_set_status, handle, "app_indicator_set_status")
	purego.RegisterLibFunc(&app_indicator_set_attention_icon_full, handle, "app_indicator_set_attention_icon_full")
	purego.RegisterLibFunc(&app_indicator_set_menu, handle, "app_indicator_set_menu")
	purego.RegisterLibFunc(&app_indicator_set_icon_full, handle, "app_indicator_set_icon_full")
	purego.RegisterLibFunc(&app_indicator_set_label, handle, "app_indicator_set_label")
	purego.RegisterLibFunc(&app_indicator_set_icon_theme_path, handle, "app_indicator_set_icon_theme_path")
	purego.RegisterLibFunc(&app_indicator_set_ordering_index, handle, "app_indicator_set_ordering_index")
	purego.RegisterLibFunc(&app_indicator_set_secondary_activate_target, handle, "app_indicator_set_secondary_activate_target")
	purego.RegisterLibFunc(&app_indicator_set_title, handle, "app_indicator_set_title")
	purego.RegisterLibFunc(&app_indicator_get_category, handle, "app_indicator_get_category")
	purego.RegisterLibFunc(&app_indicator_get_status, handle, "app_indicator_get_status")
	purego.RegisterLibFunc(&app_indicator_get_ordering_index, handle, "app_indicator_get_ordering_index")
	purego.RegisterLibFunc(&app_indicator_build_menu_from_desktop, handle, "app_indicator_build_menu_from_desktop")
}
