package codegen

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"strconv"

	"github.com/fatih/camelcase"
	"github.com/juju/errors"
)

type (
	Import struct {
		Alias string
		Path  string
	}

	Imports map[string]Import

	Field struct {
		Order int
		Name  string
		Type  string
		Tag   string
	}

	Fields map[string]Field

	Struct struct {
		Imports Imports
		Name    string
		Fields  Fields
	}

	Structs map[string]Struct

	Func struct {
		PkgPath   string
		Name      string
		Method    []string
		Protected bool
		Pattern   string
		Query     []string
		In        *Struct
		Out       *Struct
	}

	Funcs []Func

	Service struct {
		Import Import
		Funcs  Funcs
	}

	Services []Service

	ByOrder []Field
)

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

func (i Import) GetPathAsAlias() string {
	r := strings.NewReplacer("/", "", ".", "")
	return strings.ToLower(r.Replace(i.Path))
}

func (f Func) GetFullName() string {
	return f.PkgPath + ":" + f.Name
}

func (f Func) GetPkgPath() string {
	return f.PkgPath
}

func (f Func) GetOnlyPkgName() string {
	return filepath.Base(f.PkgPath)
}

func (f Func) GetMethod() []string {
	return f.Method
}

func (f Func) GetQuery() []string {
	return f.Query
}

func (f Func) GetPattern() string {
	return f.Pattern
}

func (f Func) GetHandlerName() string {
	return f.Name + "Handler"
}
func (f Func) GetProtected() bool {
	return f.Protected
}

func (fld Field) GetVarName() string {
	return makeFirstUpperCase(fld.Name)
}

func (fld Field) GetVarValue(data map[string]string) (interface{}, error) {
	// if value exists in variables
	if val, ok := data[fld.Name]; ok {
		switch fld.Type {
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
			return nil, errors.Trace(errors.Errorf("Value is complext type!Expected only primitives!"))
		}
	}

	return nil, nil
}

func (fld Field) GetVarType() string {
	return fld.Type
}

func (fld Field) GetTag() string {
	return "`json:\"" + makeTagName(fld.Name) + ",omitempty\"`"
}

func (fld Field) GetOutVarName() string {
	return makeTagName(fld.Name)
}

func (s Struct) GetName() string {
	return s.Name
}

func (s Struct) GetFieldsSlice() []Field {
	fs := []Field{}
	for _, f := range s.Fields {
		fs = append(fs, f)
	}

	sort.Sort(ByOrder(fs))

	return fs
}

func makeTagName(s string) string {
	a := []rune(s)

	if len(a) <= 3 {
		return strings.ToLower(string(a))
	}

	spl := camelcase.Split(string(a))

	if len(spl) > 1 {
		fixed := []string{}
		for _, v := range spl {
			fixed = append(fixed, makeFirstUpperCase(strings.ToLower(v)))
		}
		x := makeFirstLowerCase(strings.Join(fixed, ""))
		return x
	}

	return makeFirstLowerCase(strings.Title(string(a)))
}

func makeFirstLowerCase(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func makeFirstUpperCase(s string) string {
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
