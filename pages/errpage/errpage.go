package errpage

import (
	"html/template"
	"net/http"

	"github.com/diamondburned/l4d2lb/pages"
)

var tmpl *template.Template

func init() {
	tmpl = pages.Template("errpage", nil)
}

func RenderPage(w http.ResponseWriter, err error) {
	pages.Execute(tmpl, w, err.Error())
}

func RenderError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	RenderPage(w, err)
}
