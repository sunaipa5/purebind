package main

import (
	"strings"
	"unicode"
)

func generateParam(cType string) string {
	cType = strings.TrimSpace(cType)

	removeWords := []string{"const", "volatile", "restrict"}
	for _, w := range removeWords {
		cType = strings.ReplaceAll(cType, w, "")
	}

	cType = strings.TrimSpace(cType)
	fields := strings.Fields(cType)
	cType = strings.Join(fields, " ")
	cType = strings.ReplaceAll(cType, "*", "")

	return cType
}

func generateGoType(cType string) string {
	cType = strings.TrimSpace(cType)
	switch cType {
	case "int":
		return "int32"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "void", "void*":
		return "unsafe.Pointer"
	case "const char*", "char*":
		return "unsafe.Pointer"
	default:
		if strings.HasSuffix(cType, "*") {
			return "unsafe.Pointer"
		}
		return "unsafe.Pointer"
	}
}

func generateBindingFunc(name string) string {
	sb := strings.Builder{}
	for i, r := range name {
		if unicode.IsLetter(r) || (i > 0 && (unicode.IsDigit(r) || r == '_')) {
			sb.WriteRune(r)
		}
	}
	out := sb.String()
	if out == "" || !unicode.IsLetter(rune(out[0])) {
		out = "Func" + out
	}
	return out
}

func generateWrapperFunc(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
