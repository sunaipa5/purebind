package main

import (
	"os"
	"path"
	"regexp"
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
	"const char":     "string",
	"float":          "float32",
	"double":         "float64",
	"void":           "any",
	"size_t":         "unsafe.Pointer",
	"int64_t":        "int64",
	"uint8_t":        "uint8",
	// Gerekirse ekle
}

func convertCTypeToGo(cType string) string {
	cType = strings.TrimSpace(cType)

	funcPtrRegex := regexp.MustCompile(`^\s*([\w\s\*]+)\(\s*\*\s*(\w+)\s*\)\s*\((.*)\)\s*$`)
	if matches := funcPtrRegex.FindStringSubmatch(cType); matches != nil {
		return convertFuncPointerToGo(matches[1], matches[3])
	}

	// Pointer sayısını hesapla (sadece sondaki * leri sayar)
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
		// Parametreleri parçala virgülden
		paramsSplit := splitParams(paramList)
		for i, p := range paramsSplit {
			// Parametre tipi olarak sadece tip yazıyoruz (isimler opsiyonel)
			// Basitçe ilk kelimeleri alabiliriz
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

// Basitçe fonksiyon parametrelerini virgüle göre ayırır, parantez ve pointerleri göz ardı eder
func splitParams(params string) []string {
	var result []string
	current := strings.Builder{}
	level := 0 // parantez seviyesi
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
	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func extractFields(structBody string) string {
	lines := strings.Split(structBody, "\n")
	var sb strings.Builder

	fieldRegex := regexp.MustCompile(`^(?:/\*\*([^*]*?)\*/)?\s*(.+?)\s+(\w+)\s*;`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := fieldRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		comment := strings.TrimSpace(matches[1])
		cType := matches[2]
		name := matches[3]

		goType := convertCTypeToGo(cType)
		if comment != "" {
			sb.WriteString("\t// ")
			sb.WriteString(comment)
			sb.WriteByte('\n')
		}

		sb.WriteByte('\t')
		sb.WriteString(exportName(name))
		sb.WriteByte(' ')
		sb.WriteString(goType)
		sb.WriteByte('\n')
	}

	return sb.String()
}

func convertStruct(body string, name string) string {
	var sb strings.Builder
	sb.WriteString("type ")
	sb.WriteString(name)
	sb.WriteString(" struct {\n")
	sb.WriteString(extractFields(body))
	sb.WriteString("}\n")
	return sb.String()
}

func parseStructs(libLocation, packageName string) error {
	contentBytes, err := os.ReadFile(libLocation)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	re := regexp.MustCompile(`typedef\s+struct\s*(\w*)\s*\{([^}]*)\}\s*(\w+);`)

	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString("package ")
	sb.WriteString(packageName)
	sb.WriteString("\n\nimport \"unsafe\"\n\n")

	for _, m := range matches {
		body := m[2]
		name := m[3]
		sb.WriteString(convertStruct(body, name))
		sb.WriteByte('\n')
	}

	outputPath := path.Join(packageName, "structs.go")
	if err := os.MkdirAll(packageName, 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
