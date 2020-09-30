package errbox

import (
	"bytes"
	"html/template"

	"github.com/diamondburned/l4d2lb/pages"
)

var tmpl *template.Template

func init() {
	tmpl = pages.NewTemplate("components/errbox", nil)
}

func RenderHTML(err error) template.HTML {
	var buf bytes.Buffer
	pages.Execute(tmpl, &buf, err.Error())
	return template.HTML(buf.String())
}
