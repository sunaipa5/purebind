
package ayatana

import "unsafe"

func AppIndicatorSetStatus(self, status unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_status(self, status)
}

func AppIndicatorSetAttentionIconFull(self, icon_name, icon_desc unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_attention_icon_full(self, icon_name, icon_desc)
}

func AppIndicatorSetMenu(self, menu unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_menu(self, menu)
}

func AppIndicatorSetIconFull(self, icon_name, icon_desc unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_icon_full(self, icon_name, icon_desc)
}

func AppIndicatorSetLabel(self, label, guide unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_label(self, label, guide)
}

func AppIndicatorSetIconThemePath(self, icon_theme_path unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_icon_theme_path(self, icon_theme_path)
}

func AppIndicatorSetOrderingIndex(self, ordering_index unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_ordering_index(self, ordering_index)
}

func AppIndicatorSetSecondaryActivateTarget(self, menuitem unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_secondary_activate_target(self, menuitem)
}

func AppIndicatorSetTitle(self, title unsafe.Pointer) unsafe.Pointer {
return app_indicator_set_title(self, title)
}

func AppIndicatorGetCategory(self unsafe.Pointer) unsafe.Pointer {
return app_indicator_get_category(self)
}

func AppIndicatorGetStatus(self unsafe.Pointer) unsafe.Pointer {
return app_indicator_get_status(self)
}

func AppIndicatorGetOrderingIndex(self unsafe.Pointer) unsafe.Pointer {
return app_indicator_get_ordering_index(self)
}

func AppIndicatorBuildMenuFromDesktop(self, desktop_file, desktop_profile unsafe.Pointer) unsafe.Pointer {
return app_indicator_build_menu_from_desktop(self, desktop_file, desktop_profile)
}

