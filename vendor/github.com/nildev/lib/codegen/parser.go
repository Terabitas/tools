package codegen

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/fatih/camelcase"
)

const (
	PREFIX_INPUT  = "Input"
	PREFIX_OUTPUT = "Output"

	TokenMethod    = "@method (.*)"
	TokenPath      = "@path (.*)"
	TokenQuery     = "@query (.*)"
	TokenProtected = "@protected"

	TokenNameMethod    = "method"
	TokenNamePath      = "path"
	TokenNameQuery     = "query"
	TokenNameProtected = "protected"

	MethodPOST = "POST"
	MethodGET  = "GET"
)

var (
	availableTokens = map[string]string{
		TokenNameMethod:    TokenMethod,
		TokenNamePath:      TokenPath,
		TokenNameQuery:     TokenQuery,
		TokenNameProtected: TokenProtected,
	}
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

	f := &Func{
		PkgPath: packagePath,
		Name:    fn.Name.Name,
	}

	fnParsedComments := map[string][]string{}
	if fn.Doc != nil {
		fnParsedComments = parseAllComments(fn.Doc.List)
	}

	f.Protected = MakeProtected(fn, fnParsedComments)
	f.Method = MakeMethod(fn, fnParsedComments)
	f.Pattern = MakePattern(fn, fnParsedComments)
	f.Query = MakeQuery(fn, fnParsedComments, f.Method, f.Pattern)

	f.In, f.Out = MakeInputOutputStructs(fn, imps, comments)

	return f
}

func MakeProtected(fn *ast.FuncDecl, comments map[string][]string) bool {
	protected := false

	if v, ok := comments[TokenNameProtected]; ok {
		if len(v) > 0 {
			protected = true
		}
	}

	return protected
}

func MakeQuery(fn *ast.FuncDecl, comments map[string][]string, method, path string) []string {
	query := []string{}

	if v, ok := comments[TokenNameQuery]; ok {
		if len(v) > 1 {
			parts := strings.Split(v[1], " ")
			for _, q := range parts {
				r, _ := regexp.Compile("{(?P<name>\\w+)(?:\\s*:(?P<validation>.+))?}")
				matches := r.FindAllStringSubmatch(q, -1)
				if len(matches) > 0 {
					val := ""
					if len(matches[0]) == 3 && matches[0][2] != "" {
						val = q
					}
					query = append(query, matches[0][1], val)
				}

			}
		}
	}

	// Add these only if method is GET
	if method == MethodGET {
		// add to query params from function signature which that has no validation rules
		// and though are not added to @query tag
		if fn.Type.Params != nil {
			for _, f := range fn.Type.Params.List {

				// check if param from signature is not in path
				r, _ := regexp.Compile("{(?P<name>\\w+)(?:\\s*:(?P<validation>[^/]+))?}")
				matches := r.FindAllStringSubmatch(path, -1)
				found := false
				for _, m := range matches {
					// m [{xxx:[a-z]{4}} xxx [a-z]{4}]
					// m [{xxx} xxx]
					if len(m) > 1 {
						if m[1] == f.Names[0].Name {
							found = true
							break
						}
					}
				}

				if found {
					continue
				}

				// If it has been added while adding from @query - skip it
				// f.Names[0].Name here is `xxx` func MyF(xxx string)
				if stringInSlice(f.Names[0].Name, query) {
					continue
				}

				// if param is pointer `func MyFunc(param *string)` then it is optional - do not add it
				if _, ok := f.Type.(*ast.StarExpr); !ok {
					// if all checks passed then add it with empty string as value
					// gorilla will match any string in this case, but in query f.Names[0].Name will have to be
					// presented
					query = append(query, f.Names[0].Name, "")
				}

			}
		}
	}

	return query
}

func MakePattern(fn *ast.FuncDecl, comments map[string][]string) string {
	path := fn.Name.Name

	if v, ok := comments[TokenNamePath]; ok {
		if len(v) == 2 {
			path = v[1]
		}
	} else {
		path = "/" + makeDashedFromCamelCase(path)
	}

	return path
}

func MakeMethod(fn *ast.FuncDecl, comments map[string][]string) string {
	method := MethodGET

	if v, ok := comments[TokenNameMethod]; ok {
		if len(v) == 2 {
			method = v[1]
		}
	}

	return method
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

func parseAllComments(comments []*ast.Comment) map[string][]string {
	rez := map[string][]string{}
	if len(comments) > 0 {
		for i := 0; i < len(comments); i++ {
			for token, regPattern := range availableTokens {
				res, rezParse := parseComment(regPattern, comments[i].Text)
				if rezParse {
					rez[token] = res
				}
			}
		}
	}
	return rez
}

func makeDashedFromCamelCase(in string) string {
	splitted := camelcase.Split(in)
	return strings.ToLower(strings.Join(splitted, "-"))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
