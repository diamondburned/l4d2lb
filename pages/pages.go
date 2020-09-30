package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/diamondburned/l4d2lb/stats"
	"github.com/phogolabs/parcello"
)

//go:generate go run github.com/phogolabs/parcello/cmd/parcello -r -i *.go

func MountStatic() http.Handler {
	d, err := parcello.Manager.Dir("static/")
	if err != nil {
		log.Fatalln("Failed to access embedded `static/' directory:", err)
	}

	return http.FileServer(d)
}

func openFile(name string) (base string, content []byte) {
	// Derive a filepath using set rules.
	var fileBase = filepath.Base(name)
	var filePath = ""

	if ext := filepath.Ext(fileBase); ext != "" {
		filePath = name
		fileBase = fileBase[:len(fileBase)-len(ext)] // trim extension
	} else {
		filePath = filepath.Join(name, fileBase+".html")
	}

	f, err := parcello.Open(filePath)
	if err != nil {
		log.Fatalln("Failed to open file:", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("Failed to read file:", err)
	}

	return fileBase, b
}

func NewTemplate(name string, funcMap template.FuncMap) *template.Template {
	fileBase, b := openFile(name)

	tmpl := template.New(fileBase)
	tmpl = tmpl.Funcs(funcMap)
	tmpl = template.Must(tmpl.Parse(string(b)))

	return tmpl
}

func Execute(tmpl *template.Template, w io.Writer, v interface{}) {
	if err := tmpl.Execute(w, v); err != nil {
		log.Println("Failed to execute template:", err)
	}
}

func ExecuteComponent(tmpl *template.Template, w io.Writer, name string, v interface{}) {
	if err := tmpl.ExecuteTemplate(w, name, v); err != nil {
		log.Println("Failed to execute template:", err)
	}
}

type RenderState struct {
	*stats.Database
	SiteName string
}

func (s *RenderState) Template(name string, funcMap template.FuncMap) *Template {
	return &Template{
		Template:    NewTemplate(name, funcMap),
		RenderState: s,
	}
}

type Template struct {
	*template.Template
	*RenderState
}

func (s *Template) Execute(w io.Writer, v interface{}) {
	Execute(s.Template, w, v)
}

func (s *Template) ExecuteComponent(w io.Writer, name string, v interface{}) {
	ExecuteComponent(s.Template, w, name, v)
}

func (s *Template) RenderHTML(v interface{}) template.HTML {
	var buf bytes.Buffer
	s.Execute(&buf, v)
	return template.HTML(buf.String())
}

func (s *Template) RenderHTMLComponent(name string, v interface{}) template.HTML {
	var buf bytes.Buffer
	s.ExecuteComponent(&buf, name, v)
	return template.HTML(buf.String())
}

func (s *Template) AddComponent(path string) {
	fileBase, b := openFile(path)

	s.Template = template.Must(s.Template.Parse(fmt.Sprintf(
		"{{ define \"%s\" }}\n%s\n{{ end }}",
		fileBase, string(b),
	)))
}
