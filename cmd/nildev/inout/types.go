package inout

type (
	// Generator is responsible for generating required integration code
	// for nildev service container
	Generator interface {
		Generate(pathToServiceDir, tplPath string)
	}
)

var (
	DefaultSimpleTemplate = `package {{.PackageName}}

// THIS IS AUTO GENERATED FILE - DO NOT EDIT!
import (
	{{range .Imports}}
		{{.Alias}} "{{.Path}}"
	{{end}}
)

type (
	{{range $i, $f := .Funcs}}

	// {{$f.In.GetName}} struct
	{{$f.In.GetName}} struct {
		{{range $j, $e := $f.In.GetFieldsSlice}}
		{{$e.GetVarName}} {{$e.GetVarType}} {{$e.GetTag}}
		{{end}}
	}

	// {{$f.Out.GetName}} struct
	{{$f.Out.GetName}} struct {
		{{range $j, $e := $f.Out.GetFieldsSlice}}
		{{$e.GetVarName}} {{$e.GetVarType}} {{$e.GetTag}}
		{{end}}
	}
	{{end}}
)

{{range $i, $f := .Funcs}}
// {{$f.GetHandlerName}} HTTP request handler
func {{$f.GetHandlerName}}(rw http.ResponseWriter, r *http.Request) {
	returnCode := http.StatusOK

	reqDTO := &{{$f.In.Name}}{}
	utils.UnmarshalRequest(r, reqDTO)

	{{range $j, $e := $f.Out.GetFieldsSlice}}{{if $j}},{{end}} {{$e.GetOutVarName}}{{end}} := {{$f.Name}}({{range $j, $e := $f.In.GetFieldsSlice}}{{if $j}},{{end}}reqDTO.{{$e.GetVarName}}{{end}})

	if err != nil {
		returnCode = http.StatusInternalServerError
	}

	respDTO := &{{$f.Out.Name}}{
		{{range $j, $e := $f.Out.GetFieldsSlice}}
		{{$e.GetVarName}}:{{$e.GetOutVarName}},
		{{end}}
	}

	utils.Respond(rw, respDTO, returnCode)
}
{{end}}

// NildevRoutes returns routes to be registered
func NildevRoutes() router.Routes {
	routes := router.Routes{
		BasePattern: "{{.BasePattern}}",
		Routes: make([]router.Route, {{.RoutesNum}}),
	}

	{{range $i, $f := .Funcs}}
	routes.Routes[{{$i}}] = router.Route{
		Name:        "{{$f.GetFullName}}",
		Method:      "{{$f.GetMethod}}",
		Pattern:     "{{$f.GetPattern}}",
		HandlerFunc: {{$f.GetHandlerName}},
	}
	{{end}}

	return routes
}
`
)
