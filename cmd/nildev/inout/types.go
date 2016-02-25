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
    httpCodeSetByUser := false
    returnCode := http.StatusOK

    {{if $f.In.GetFieldsSlice}}
    var requestData map[string]string
    requestData = mux.Vars(r)
    log.Infof("Request data [%+v]",requestData)
    {{end}}

    reqDTO := &{{$f.In.Name}}{}
    utils.UnmarshalRequest(r, reqDTO)

    {{range $j, $e := $f.In.GetFieldsSlice}}
    {{if eq $e.Name "user"}}

    user := context.Get(r, "user")
    if user != nil {
        reqDTO.{{$e.GetVarName}} = registry.MakeUser(user.(*jwt.Token).Claims)
    }

    {{else}}

    cv{{$e.Name}}, convErr := GetVarValue(requestData, "{{$e.Name}}", "{{$e.Type}}")
    if convErr != nil {
	returnCode = http.StatusInternalServerError
	utils.Respond(rw, convErr, returnCode)
	return
    }

    if cv{{$e.Name}} != nil  {
	reqDTO.{{$e.GetVarName}} = cv{{$e.Name}}.({{$e.Type}})
    }

    {{end}}
    {{end}}

    {{range $j, $e := $f.Out.GetFieldsSlice}}{{if $j}},{{end}} {{$e.GetOutVarName}}{{end}} := {{$f.Name}}({{range $j, $e := $f.In.GetFieldsSlice}}{{if $j}},{{end}}reqDTO.{{$e.GetVarName}}{{end}})

    {{range $j, $e := $f.Out.GetFieldsSlice}}
    {{if eq $e.Name "httpHeaders"}}

    hh := {{$e.GetOutVarName}}
    for k, v := range hh {
    	if len(v) == 1 {
    	    rw.Header().Set(k, v[0])
    	} else if len(v) > 1{
	    for _, vv := range v {
	        rw.Header().Add(k, vv)
	    }
    	}
    }
    {{end}}

    {{if eq $e.Name "httpStatus"}}
    httpCodeSetByUser = true
    returnCode = {{$e.GetOutVarName}}
    {{end}}

    {{end}}

    if !httpCodeSetByUser {
        if err != nil {
            returnCode = http.StatusInternalServerError
	}
    }

    respDTO := &{{$f.Out.Name}}{
	{{range $j, $e := $f.Out.GetFieldsSlice}}
		{{if ne $e.Name "httpHeaders"}}
			{{if ne $e.Name "httpStatus"}}
				{{$e.GetVarName}}:{{$e.GetOutVarName}},
			{{end}}
		{{end}}
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
		Protected:   {{$f.GetProtected}},
		HandlerFunc: {{$f.GetHandlerName}},
		Queries:     []string{
		    {{range $f.GetQuery}}
		    "{{.}}",
		    {{end}}
		},
	}
	{{end}}

	return routes
}

func GetVarValue(data map[string]string, name, typ string) (interface{}, error) {
	// if value exists in variables
	if val, ok := data[name]; ok {
		switch typ {
		case "string":
			return val, nil
		case "*string":
			return &val, nil
		case "*int":
			i, err := strconv.Atoi(val)
			return &i, err
		case "int":
			i, err := strconv.Atoi(val)
			return i, err
		case "int8":
			i, err := strconv.ParseInt(val, 10, 8)
			return int8(i), err
		case "*int8":
			i, err := strconv.ParseInt(val, 10, 8)
			icst := int8(i)
			return &icst, err
		case "int16":
			i, err := strconv.ParseInt(val, 10, 16)
			return int16(i), err
		case "*int16":
			i, err := strconv.ParseInt(val, 10, 16)
			icst := int16(i)
			return &icst, err
		case "int32":
			i, err := strconv.ParseInt(val, 10, 32)
			return int32(i), err
		case "*int32":
			i, err := strconv.ParseInt(val, 10, 32)
			icst := int32(i)
			return &icst, err
		case "int64":
			i, err := strconv.ParseInt(val, 10, 64)
			return i, err
		case "*int64":
			i, err := strconv.ParseInt(val, 10, 64)
			return &i, err
		case "float32":
			i, err := strconv.ParseFloat(val, 32)
			return float32(i), err
		case "*float32":
			i, err := strconv.ParseFloat(val, 32)
			icst := float32(i)
			return &icst, err
		case "float64":
			i, err := strconv.ParseFloat(val, 64)
			return i, err
		case "*float64":
			i, err := strconv.ParseFloat(val, 64)
			icst := float64(i)
			return &icst, err
		case "bool":
			b, err := strconv.ParseBool(val)
			return b, err
		case "*bool":
			b, err := strconv.ParseBool(val)
			return &b, err
		default:
			return nil, nil
		}
	}

	return nil, nil
}
`
)
