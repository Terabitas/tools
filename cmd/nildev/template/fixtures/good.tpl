//nildev:template nildev:test-good v0.1.9
package {{.PackageName}}

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
