package main

import (
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var knownCTypes = map[string]string{
	"int":            "int",
	"unsigned int":   "uint",
	"short":          "int16",
	"unsigned short": "uint16",
	"long":           "int64",
	"unsigned long":  "uint64",
	"char":           "byte",
	"signed char":    "int8",
	"unsigned char":  "byte",
	"const char":     "uintptr",
	"float":          "float32",
	"double":         "float64",
	"void":           "byte",
	"size_t":         "uintptr",
	"int64_t":        "int64",
	"uint8_t":        "uint8",
}

var goReservedWords = map[string]bool{
	"break": true, "default": true, "func": true, "interface": true, "select": true,
	"case": true, "defer": true, "go": true, "map": true, "struct": true,
	"chan": true, "else": true, "goto": true, "package": true, "switch": true,
	"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
	"continue": true, "for": true, "import": true, "return": true, "var": true,
}

func convertCTypeToGo(cType string) string {
	cType = strings.TrimSpace(cType)

	funcPtrRegex := regexp.MustCompile(`^\s*([\w\s\*]+)\(\s*\*\s*(\w+)\s*\)\s*\((.*)\)\s*$`)
	if matches := funcPtrRegex.FindStringSubmatch(cType); matches != nil {
		return convertFuncPointerToGo(matches[1], matches[3])
	}

	pointerCount := 0
	for strings.HasSuffix(cType, "*") {
		pointerCount++
		cType = strings.TrimSpace(cType[:len(cType)-1])
	}

	for {
		switch {
		case strings.HasPrefix(cType, "const "):
			cType = strings.TrimSpace(cType[len("const "):])
		case strings.HasPrefix(cType, "struct "):
			cType = strings.TrimSpace(cType[len("struct "):])
		case strings.HasPrefix(cType, "enum "):
			cType = strings.TrimSpace(cType[len("enum "):])
		default:
			goto breakLoop
		}
	}

breakLoop:

	goType, found := knownCTypes[cType]
	if !found {
		goType = cType
	}

	if goType != "string" {
		goType = strings.Repeat("*", pointerCount) + goType
	} else if pointerCount > 0 {
		goType = "*" + goType
	}

	return goType
}

func convertFuncPointerToGo(returnType, params string) string {
	goReturnType := convertCTypeToGo(returnType)

	paramList := strings.TrimSpace(params)
	if paramList == "void" || paramList == "" {
		paramList = ""
	} else {
		paramsSplit := splitParams(paramList)
		for i, p := range paramsSplit {
			parts := strings.Fields(p)
			if len(parts) == 0 {
				continue
			}
			paramsSplit[i] = convertCTypeToGo(parts[0])
		}
		paramList = strings.Join(paramsSplit, ", ")
	}

	return "func(" + paramList + ") " + goReturnType
}

func splitParams(params string) []string {
	var result []string
	current := strings.Builder{}
	level := 0
	for _, r := range params {
		switch r {
		case '(':
			level++
			current.WriteRune(r)
		case ')':
			level--
			current.WriteRune(r)
		case ',':
			if level == 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}
	return result
}

func exportName(name string) string {
	if name == "" {
		return ""
	}

	tempName := name
	if goReservedWords[name] {
		tempName = name + "Field"
	}

	runes := []rune(tempName)
	runes[0] = unicode.ToUpper(runes[0])
	exported := string(runes)

	switch exported {
	case "Int", "String", "Float", "Bool", "Byte", "Uintptr":
		return exported + "_"
	}

	return exported
}

func extractFields(structBody string) string {
	lines := strings.Split(structBody, "\n")
	var sb strings.Builder

	fieldRegex := regexp.MustCompile(`^(?:/\*\*([^*]*?)\*/)?\s*(.+?[\s\*]+)(\w+)((?:\[\d+\])+)?\s*;`)
	funcPtrRegex := regexp.MustCompile(`^(?:/\*\*([^*]*?)\*/)?\s*(.+?)\s*\(\s*\*\s*(\w+)((?:\[\d+\])+)?\s*\)\s*\((.*)\)\s*;`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if matches := fieldRegex.FindStringSubmatch(line); matches != nil {
			comment := strings.TrimSpace(matches[1])
			cType := strings.TrimSpace(matches[2])
			name := matches[3]
			arrayPart := matches[4]

			goType := convertCTypeToGo(cType)
			if arrayPart != "" {
				goType = arrayPart + goType
			}

			writeField(&sb, comment, name, goType)
			continue
		}

		if matches := funcPtrRegex.FindStringSubmatch(line); matches != nil {
			comment := strings.TrimSpace(matches[1])
			name := matches[3]
			arrayPart := matches[4]

			goType := "uintptr"
			if arrayPart != "" {
				goType = arrayPart + goType
			}

			writeField(&sb, comment, name, goType)
			continue
		}
	}
	return sb.String()
}

func writeField(sb *strings.Builder, comment, name, goType string) {
	if comment != "" {
		sb.WriteString("\t// " + comment + "\n")
	}
	sb.WriteString("\t" + exportName(name) + " " + goType + "\n")
}

func convertStruct(body string, name string) string {
	var sb strings.Builder
	sb.WriteString("type ")
	sb.WriteString(name)
	sb.WriteString(" struct {\n")

	fields := extractFields(body)

	if fields == "" {
		sb.WriteString("\t_ uintptr // Padding for empty C struct\n")
	} else {
		sb.WriteString(fields)
	}

	sb.WriteString("}\n")
	return sb.String()
}

func parseStructs(libLocation, packageName string) error {
	contentBytes, err := os.ReadFile(libLocation)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	var sb strings.Builder
	sb.WriteString("package " + packageName + "\n\n")

	// ------------- ENUMS START -------------
	reEnum := regexp.MustCompile(`(?s)typedef\s+enum\s*\{([^}]*)\}\s*(\w+);`)
	reEnumItem := regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*(?:=\s*([^,//\s]+))?`)
	for _, m := range reEnum.FindAllStringSubmatch(content, -1) {
		if len(m) < 3 {
			continue
		}
		body, enumName := m[1], m[2]

		var enumBody strings.Builder
		lastValue := -1
		hasValidMember := false

		lines := strings.Split(body, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
				continue
			}

			match := reEnumItem.FindStringSubmatch(line)
			if match == nil {
				continue
			}

			name := match[1]
			valPart := match[2]

			if name == "false" || name == "true" || name == "type" || name == "func" {
				continue
			}

			hasValidMember = true

			if valPart != "" {
				valPart = strings.Trim(valPart, ";,")
				enumBody.WriteString("\t" + name + " = " + valPart + "\n")

				if i, err := strconv.ParseInt(valPart, 0, 64); err == nil {
					lastValue = int(i)
				}
			} else {
				lastValue++
				enumBody.WriteString("\t" + name + " = " + strconv.Itoa(lastValue) + "\n")
			}
		}

		if hasValidMember {
			sb.WriteString("// " + enumName + "\nconst (\n")
			sb.WriteString(enumBody.String())
			sb.WriteString(")\n\n")
		}
	}

	//------------- Alias -------------
	reAlias := regexp.MustCompile(`(?m)^\s*typedef\s+([\w\s\*]+)\s+(\w+);`)
	for _, m := range reAlias.FindAllStringSubmatch(content, -1) {
		if len(m) < 3 {
			continue
		}
		fullOldType := strings.TrimSpace(m[1])
		newName := strings.TrimSpace(m[2])
		if strings.HasPrefix(fullOldType, "struct") || strings.HasPrefix(fullOldType, "enum") {
			continue
		}

		sb.WriteString("type " + newName + " = " + convertCTypeToGo(fullOldType) + "\n")
	}
	sb.WriteString("\n")

	//------------- Structs -------------
	reStruct := regexp.MustCompile(`(?s)typedef\s+struct\s*\w*\s*\{([^}]*)\}\s*(\w+);`)
	for _, m := range reStruct.FindAllStringSubmatch(content, -1) {
		if len(m) < 3 {
			continue
		}
		sb.WriteString(convertStruct(m[1], m[2]) + "\n")
	}

	return os.WriteFile(path.Join(packageName, "structs.go"), []byte(sb.String()), 0644)
}
