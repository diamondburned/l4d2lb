package errpage

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/diamondburned/l4d2lb/pages"
)

var tmpl *template.Template

func init() {
	tmpl = pages.NewTemplate("errpage", template.FuncMap{
		"shorten": func(err string) string {
			// Grab the first sentence.
			first := strings.TrimSpace(strings.Split(err, ":")[0]) + "."

			// Capitalize first letter.
			r, sz := utf8.DecodeRuneInString(first)
			first = string(unicode.ToUpper(r)) + first[sz:]

			return first
		},
	})
}

func RenderPage(w http.ResponseWriter, err error) {
	pages.Execute(tmpl, w, err.Error())
}

func RenderError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	RenderPage(w, err)
}

func RenderHTML(err error) template.HTML {
	var buf bytes.Buffer
	pages.Execute(tmpl, &buf, err.Error())
	return template.HTML(buf.String())
}
