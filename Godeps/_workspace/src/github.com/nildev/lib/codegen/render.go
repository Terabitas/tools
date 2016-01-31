package codegen

import (
	"io"
	"text/template"
)

// Render method
func Render(w io.Writer, tpl string, data interface{}) error {
	templateFuncs := template.FuncMap{}
	tplParsed := template.Must(template.New("render").Funcs(templateFuncs).Parse(tpl))
	return tplParsed.Execute(w, data)
}
