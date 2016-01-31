package routes

type (
	// Generator is responsible for generating required integration code
	// for nildev service container
	Generator interface {
		Generate(pathToServiceContainerDir string, pathToServices []string, tplPath string)
	}
)

var (
	DefaultTemplate = `package {{.PackageName}}
import (
	{{range .Imports}}
		{{.Alias}} "{{.Path}}"
	{{end}}
	{{range .Services}}
		{{.Import.GetPathAsAlias}} "{{.Import.Path}}"
	{{end}}
)

func BuildRoutes() []router.Routes {
	routes := make([]router.Routes, {{.RoutesNum}})

	{{range $i, $f := .Services}}
	routes = append(routes, {{.Import.GetPathAsAlias}}.NildevRoutes())
	{{end}}

	return routes
}
`
)
