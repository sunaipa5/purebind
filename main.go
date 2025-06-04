package main

import (
	"embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

func werror(msg string, err error) {
	fmt.Printf("\033[31m[ ERROR ][ %s ][ %v ] \033[0m\n", msg, err)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <project-name> <header-file-path>")
		return
	}

	projectName := os.Args[1]
	libLocation := os.Args[2]

	funcs, err := parseHeader(libLocation)
	if err != nil {
		werror("Failed to find header file", err)
		return
	}

	data := TemplateData{
		PackageName: projectName,
		Functions:   funcs,
	}

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

	//Wrapper.t
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

	fmt.Println("Done.")
}
