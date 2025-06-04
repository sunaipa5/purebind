package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Param struct {
	Name string
	Type string
}

type Function struct {
	CName         string
	CFunc         string
	GOFunc        string
	ReturnType    string
	ParamTypes    string
	ParamNames    string
	WrapperParams string
}

type TemplateData struct {
	PackageName string
	Functions   []Function
}

func parseHeader(filename string) ([]Function, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`(?m)^([a-zA-Z_][a-zA-Z0-9_\*\s]+)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\);`)
	matches := reg.FindAllStringSubmatch(string(content), -1)

	var funcs []Function
	for _, match := range matches {
		retType := strings.TrimSpace(match[1])
		name := match[2]
		args := strings.TrimSpace(match[3])

		var params []Param
		if args != "void" && args != "" {
			argList := strings.Split(args, ",")

			for i, arg := range argList {
				arg = strings.TrimSpace(arg)
				parts := strings.Fields(arg)
				if len(parts) == 0 {
					continue
				}

				var argType, paramName string
				if len(parts) == 1 {
					argType = parts[0]
					paramName = fmt.Sprintf("p%d", i)
				} else {
					paramName = parts[len(parts)-1]
					argType = strings.Join(parts[:len(parts)-1], " ")
				}
				params = append(params, Param{
					Name: generateParam(paramName),
					Type: generateGoType(argType),
				})
			}

		}

		//Group parameters by data type
		groupedParams := make(map[string][]string)
		for _, p := range params {
			groupedParams[p.Type] = append(groupedParams[p.Type], p.Name)
		}

		// Binding function parameters - definition (types)
		var funcParamTypes []string
		for _, p := range params {
			funcParamTypes = append(funcParamTypes, p.Type)
		}

		// Binding function parameters - usage (names)
		var funcParamNames []string
		for _, names := range groupedParams {
			funcParamNames = append(funcParamNames, strings.Join(names, ", "))
		}

		//Wrapper function parameters
		var wrapperParams []string
		for typ, names := range groupedParams {
			wrapperParams = append(wrapperParams, strings.Join(names, ", ")+" "+typ)
		}

		funcs = append(funcs, Function{
			CName:         name,
			CFunc:         generateBindingFunc(name),
			GOFunc:        generateWrapperFunc(name),
			ParamTypes:    strings.Join(funcParamTypes, ","),
			ParamNames:    strings.Join(funcParamNames, ","),
			WrapperParams: strings.Join(wrapperParams, ", "),
			ReturnType:    generateGoType(retType),
		})
	}

	return funcs, nil
}
