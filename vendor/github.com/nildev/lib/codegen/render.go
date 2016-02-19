package codegen

import (
	"io"
	"reflect"
	"text/template"
)

// Render method
func Render(w io.Writer, tpl string, data interface{}) error {
	templateFuncs := template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
	}
	tplParsed := template.Must(template.New("render").Funcs(templateFuncs).Parse(tpl))
	return tplParsed.Execute(w, data)
}
