package s3n

import (
	"html/template"
	"path/filepath"
	"os"
)

var Template *template.Template = nil

func T() *template.Template {
	return Template
}

func Loadtemplates() {
	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}
	dir, err = filepath.Abs(dir+"/template")
	if err!=nil {
		panic(err)
	}
	
	pattern := filepath.Join(dir, "*.tmpl")
	Template = template.Must(template.ParseGlob(pattern))
}
