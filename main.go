package main

import (
	"embed"
	"fmt"
	"os"
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

	parseTemplates(data)

	fmt.Println("Done.")
}
