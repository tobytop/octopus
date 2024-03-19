package main

import (
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	projectName = flag.String("project_name", "", "project name can get in go.mod file")
	outPath     = flag.String("out_path", "/", "same as the go out path")
	templete    = `
	package proto_menu
	
	import ({{ range .Packages }}
	_ "{{ . }}"
	{{- end}}
	)
	`
)

type templetaData struct {
	Packages []string
}

func main() {
	flag.Parse()
	protogen.Options{
		ParamFunc: flag.Set,
	}.Run(func(gen *protogen.Plugin) error {
		abspath, err := filepath.Abs(*outPath)
		if err != nil {
			return err
		}
		index := strings.LastIndex(abspath, *projectName)
		imporPrefix := strings.ReplaceAll(abspath[index:], "\\", "/")
		protoPaths := make([]string, 0)
		for _, f := range gen.Files {
			path := imporPrefix + "/" + string(f.GoImportPath)
			protoPaths = append(protoPaths, path[:len(path)-1])
		}
		t, err := template.New("menu").Parse(templete)
		if err != nil {
			return err
		}

		path := *outPath + "/proto_menu"
		fileName := path + "/proto_menu.go"
		os.Remove(fileName)
		os.Remove(path)
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		err = t.Execute(f, templetaData{
			Packages: protoPaths,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
