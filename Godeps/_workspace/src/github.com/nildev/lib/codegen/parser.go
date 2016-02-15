package codegen

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

const (
	PREFIX_INPUT  = "Input"
	PREFIX_OUTPUT = "Output"

	TokenMethod = "@method (.*)"

	MethodPOST = "POST"
)

func ParsePackage(comments []*ast.CommentGroup) *string {
	if len(comments) == 0 {
		return nil
	}

	if len(comments[0].List) == 0 {
		return nil
	}

	// Canonical import path https://golang.org/doc/go1.4#canonicalimports
	impStmt := comments[0].List[0].Text
	if !strings.Contains(impStmt, "// import") {
		return nil
	}

	packagePath := strings.Trim(strings.TrimLeft(impStmt, "// import"), "\"")

	return &packagePath
}

func MakeFunc(fn *ast.FuncDecl, imps []*ast.ImportSpec, comments []*ast.CommentGroup) *Func {
	method := MethodPOST
	if len(comments) == 0 {
		return nil
	}

	if len(comments[0].List) == 0 {
		return nil
	}

	// Canonical import path https://golang.org/doc/go1.4#canonicalimports
	impStmt := comments[0].List[0].Text
	if !strings.Contains(impStmt, "// import") {
		return nil
	}

	packagePath := strings.Trim(strings.TrimLeft(impStmt, "// import"), "\"")

	if len(fn.Doc.List) > 0 {
		for i := 0; i < len(fn.Doc.List); i++ {
			res, rez := parseComment(TokenMethod, fn.Doc.List[i].Text)
			if rez {
				parts := strings.Split(res[0], " ")
				if len(parts) == 2 {
					method = parts[1]
				}
			}
		}
	}

	f := &Func{
		PkgPath: packagePath,
		Method:  method,
		Pattern: fn.Name.Name,
		Name:    fn.Name.Name,
	}

	f.In, f.Out = MakeInputOutputStructs(fn, imps, comments)

	return f
}

func MakeInputOutputStructs(fn *ast.FuncDecl, imps []*ast.ImportSpec, comments []*ast.CommentGroup) (*Struct, *Struct) {
	var input *Struct
	if fn.Type.Params != nil {
		input = &Struct{
			Name:    makeInputStructName(fn.Name.Name),
			Fields:  makeFields(fn.Type.Params),
			Imports: makeImports(imps),
		}
	}

	var output *Struct
	if fn.Type.Results != nil {
		output = &Struct{
			Name:    makeOutputStructName(fn.Name.Name),
			Fields:  makeFields(fn.Type.Results),
			Imports: makeImports(imps),
		}
	}

	return input, output
}

func makeInputStructName(fnName string) string {
	return PREFIX_INPUT + fnName
}

func makeOutputStructName(fnName string) string {
	return PREFIX_OUTPUT + fnName
}

func makeImports(imps []*ast.ImportSpec) Imports {
	imports := Imports{}
	for _, i := range imps {
		alias := ""
		key := strings.Trim(i.Path.Value, "\"")
		if i.Name != nil {
			key = i.Name.Name
			alias = key
		}
		imports[key] = Import{
			Alias: alias,
			Path:  strings.Trim(i.Path.Value, "\""),
		}
	}
	return imports
}

func makeFields(fields *ast.FieldList) Fields {
	fs := Fields{}
	index := 0
	for _, f := range fields.List {
		cf := makeField(f)
		cf.Order = index
		fs[cf.Name] = cf
		index++
	}

	return fs
}

func makeField(field *ast.Field) Field {
	//	fmt.Printf("[%s] %T \n", field.Names[0].Name, field.Type)
	f := Field{
		Name: field.Names[0].Name,
		Type: makeFieldType(field.Type),
	}

	return f
}

func makeFieldType(expr ast.Expr) string {
	var s string
	switch k := expr.(type) {
	case *ast.Ellipsis:
		s = makeElipsisType(k)
	case *ast.StructType:
		s = makeStructType(k)
	case *ast.StarExpr:
		s = makePtrType(k)
	case *ast.SelectorExpr:
		s = makeSelectorType(k)
	case *ast.MapType:
		s = makeMapType(k)
	case *ast.ArrayType:
		s = makeArrayType(k)
	case *ast.InterfaceType:
		s = makeInterfaceType(k)
	case *ast.Ident:
		s = fmt.Sprintf("%s", k)
	}

	return s
}

func makeInterfaceType(i *ast.InterfaceType) string {
	return fmt.Sprintf("interface{}")
}

func makeElipsisType(e *ast.Ellipsis) string {
	return fmt.Sprintf("%s", makeFieldType(e.Elt))
}

func makeMapType(m *ast.MapType) string {
	return fmt.Sprintf("map[%s]%s", makeFieldType(m.Key), makeFieldType(m.Value))
}

func makeSelectorType(sel *ast.SelectorExpr) string {
	return fmt.Sprintf("%s.%s", sel.X, sel.Sel)
}

func makeStructType(strc *ast.StructType) string {
	return fmt.Sprintf("struct{}")
}

func makePtrType(ptr *ast.StarExpr) string {
	return fmt.Sprintf("*%s", ptr.X)
}

func makeArrayType(a *ast.ArrayType) string {
	l := ""
	if a.Len != nil {
		l = fmt.Sprintf("%s", a.Len)
	}
	return fmt.Sprintf("[%s]%s", l, makeFieldType(a.Elt))
}

func parseComment(pattern string, comment string) ([]string, bool) {
	result := []string{}
	tpl, err := regexp.Compile(pattern)

	if err != nil {
		return []string{}, false
	}

	if tpl.MatchString(comment) {
		result = tpl.FindStringSubmatch(comment)

		return result, true
	}

	return result, false
}
