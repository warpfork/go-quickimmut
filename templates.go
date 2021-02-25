package main

import (
	"io"
	"text/template"
)

type mapInfo struct {
	KeyType    string // may or may not have '*' prefix, may be upper or lowercase
	ValueType  string // may or may not have '*' prefix, may be upper or lowercase
	Visibility bool   // true to make the carrier types exported
	AttachTo   string // if nonzero, constructor funcs attached to this type (which you should've declared elsewhere in your not-generated source).
}

func (mi mapInfo) Eval(w io.Writer) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			"upper": upper,
		}).Parse(mapTmpl))
	if err := tmpl.Execute(w, mi); err != nil {
		panic(err)
	}
}
func (mi mapInfo) Name() string {
	startCase := lower
	if mi.Visibility {
		startCase = upper
	}
	return startCase("map") + "__" + upper(nostar(mi.KeyType)) + "__" + upper(nostar(mi.ValueType))
}
func (mi mapInfo) Method() string {
	if mi.AttachTo != "" {
		return "(" + mi.AttachTo + ") "
	}
	return ""
}

// Templated form of https://play.golang.org/p/a-DqVuGOfmn .
//
// Writing this template has been a fascinating educational experience in being reminded how redundant golang is.
// Now, imagine how many more times "KeyType" and "ValueType" would've been repeated if I hadn't pre-templated them into "Name".
// (Maybe the forthcoming generics would make this a bit better.  No idea.)
var mapTmpl = `
type {{ .Name }} struct {
	x map[{{ .KeyType }}]{{ .ValueType }}
}
type {{ .Name }}__Entry struct {
	k {{ .KeyType }}
	v {{ .ValueType }}
}
type {{ .Name }}__Builder {{ .Name }}

func {{ .Method -}} Make{{ .Name | upper }}(ents ...{{ .Name }}__Entry) {{ .Name }} {
	x := make(map[{{ .KeyType }}]{{ .ValueType }}, len(ents))
	for _, y := range ents {
		x[y.k] = y.v
	}
	return {{ .Name }}{x}
}
func {{ .Method -}} Make{{ .Name | upper }}__Entry(k {{ .KeyType }}, v {{ .ValueType }}) {{ .Name }}__Entry {
	return {{ .Name }}__Entry{k, v}
}
func {{ .Method -}} Start{{ .Name | upper }}(sizeHint int) {{ .Name }}__Builder {
	return {{ .Name }}__Builder{make(map[{{ .KeyType }}]{{ .ValueType }}, sizeHint)}
}
func (b *{{ .Name }}__Builder) Append(k {{ .KeyType }}, v {{ .ValueType }}) {
	b.x[k] = v
}
func (b *{{ .Name }}__Builder) Finish() {{ .Name }} {
	v := *b
	b.x = nil
	return {{ .Name }}(v)
}
`

type listInfo struct {
	ValueType  string // may or may not have '*' prefix, may be upper or lowercase
	Visibility bool   // true to make the carrier types exported
	AttachTo   string // if nonzero, constructor funcs attached to this type (which you should've declared elsewhere in your not-generated source).
}

func (li listInfo) Eval(w io.Writer) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			"upper": upper,
		}).Parse(listTmpl))
	if err := tmpl.Execute(w, li); err != nil {
		panic(err)
	}
}
func (li listInfo) Name() string {
	startCase := lower
	if li.Visibility {
		startCase = upper
	}
	return startCase("list") + "__" + upper(nostar(li.ValueType))
}
func (li listInfo) Method() string {
	if li.AttachTo != "" {
		return "(" + li.AttachTo + ") "
	}
	return ""
}

// Templated form of https://play.golang.org/p/YMDuCUCrf-4 .
//
// The 'copy' is necessary because if invocation is of the `Fn(slice...)` form, Go passes the (mutable!) slice reference in.
// The 'copy' call is less problematic than one might first expect, however,
// because if the invocation is in varargs form, that slice doesn't escape (and thus we don't get multiple heap alloc costs).
// (The need for the 'copy' becomes somewhat mooted if your value type is unexported (because a `[]unexported` can't be created outside your package),
// but this generator tool accepts exported value types, and it's also darn hard to validate that no other slice references are leaked by your package,
// so it seems reasonable to do this defense unconditionally.)
var listTmpl = `
type {{ .Name }} struct {
	x []{{ .ValueType }}
}
type {{ .Name }}__Builder {{ .Name }}

func {{ .Method -}} Make{{ .Name | upper }}(ents ...{{ .ValueType }}) {{ .Name }} {
	x := make([]{{ .ValueType }}, len(ents))
	copy(x, ents)
	return {{ .Name }}{x}
}
func {{ .Method -}} Start{{ .Name | upper }}(sizeHint int) {{ .Name }}__Builder {
	return {{ .Name }}__Builder{make([]{{ .ValueType }}, 0, sizeHint)}
}
func (b *{{ .Name }}__Builder) Append(v {{ .ValueType }}) {
	b.x = append(b.x, v)
}
func (b *{{ .Name }}__Builder) Finish() {{ .Name }} {
	v := *b
	b.x = nil
	return {{ .Name }}(v)
}
`
