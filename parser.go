package main

import (
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
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
					paramName = "p" + strconv.Itoa(i)
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

		// Binding function parameters - definition (types)
		var funcParamTypes []string
		for _, p := range params {
			funcParamTypes = append(funcParamTypes, p.Type)
		}

		// Binding function parameters - usage (names)
		var funcParamNames []string
		for _, p := range params {
			funcParamNames = append(funcParamNames, p.Name)
		}

		// Wrapper function parameters
		var wrapperParams []string
		for _, p := range params {
			wrapperParams = append(wrapperParams, p.Name+" "+p.Type)
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

func parseTemplates(data TemplateData) {
	//Lib.tmpl
	tmpl, err := template.ParseFS(templatesFS, "templates/lib.tmpl")
	if err != nil {
		werror("Failed to parse lib template", err)
		return
	}
	err = os.MkdirAll(data.PackageName, 0755)
	if err != nil {
		werror("Failed to create directory", err)
		return
	}

	outFile, err := os.Create(path.Join(data.PackageName, "lib.go"))
	if err != nil {
		werror("Failed to create lib file", err)
		return
	}
	err = tmpl.ExecuteTemplate(outFile, "bindings", data)
	if err != nil {
		werror("Failed to execute lib template", err)
		return
	}
	outFile.Close()

	//Wrapper.tmpl
	wrapperTmpl, err := template.ParseFS(templatesFS, "templates/wrapper.tmpl")
	if err != nil {
		werror("Failed to parse wrapper template", err)
		return
	}
	wrapperFile, err := os.Create(path.Join(data.PackageName, "wrapper.go"))
	if err != nil {
		werror("Failed to create wrapper file", err)
		return
	}
	err = wrapperTmpl.ExecuteTemplate(wrapperFile, "bindings", data)
	if err != nil {
		werror("Failed to execute wrapper template", err)
		return
	}
	wrapperFile.Close()

	//linux.tmpl
	linuxTmpl, err := template.ParseFS(templatesFS, "templates/linux.tmpl")
	if err != nil {
		werror("Failed to parse wrapper template", err)
		return
	}

	linuxFile, err := os.Create(path.Join(data.PackageName, "linux.go"))
	if err != nil {
		werror("Failed to create wrapper file", err)
		return
	}
	err = linuxTmpl.ExecuteTemplate(linuxFile, "bindings", data)
	if err != nil {
		werror("Failed to execute wrapper template", err)
		return
	}
	linuxFile.Close()

	//windows.tmpl
	windowsTmpl, err := template.ParseFS(templatesFS, "templates/windows.tmpl")
	if err != nil {
		werror("Failed to parse wrapper template", err)
		return
	}

	windowsFile, err := os.Create(path.Join(data.PackageName, "windows.go"))
	if err != nil {
		werror("Failed to create wrapper file", err)
		return
	}
	err = windowsTmpl.ExecuteTemplate(windowsFile, "bindings", data)
	if err != nil {
		werror("Failed to execute wrapper template", err)
		return
	}
	windowsFile.Close()

	//darwin.tmpl
	darwinTmpl, err := template.ParseFS(templatesFS, "templates/darwin.tmpl")
	if err != nil {
		werror("Failed to parse wrapper template", err)
		return
	}

	darwinFile, err := os.Create(path.Join(data.PackageName, "darwin.go"))
	if err != nil {
		werror("Failed to create wrapper file", err)
		return
	}
	err = darwinTmpl.ExecuteTemplate(darwinFile, "bindings", data)
	if err != nil {
		werror("Failed to execute wrapper template", err)
		return
	}
	darwinFile.Close()
}
