package pages

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/phogolabs/parcello"
)

//go:generate go run github.com/phogolabs/parcello/cmd/parcello -r -i *.go

func Template(name string) *template.Template {
	f, err := parcello.Open(filepath.Join(name, name+".html"))
	if err != nil {
		log.Fatalln("Failed to open file:", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("Failed to read file:", err)
	}

	tmpl := template.New(name)
	tmpl = template.Must(tmpl.Parse(string(b)))
	return tmpl
}

func Execute(tmpl *template.Template, w io.Writer, v interface{}) {
	if err := tmpl.Execute(w, v); err != nil {
		log.Println("Failed to execute template:", err)
	}
}
