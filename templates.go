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
	return startCase(nostar(mi.KeyType)) + upper(nostar(mi.ValueType))
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
type {{ .Name }}Map struct {
	x map[{{ .KeyType }}]{{ .ValueType }}
}
type {{ .Name }}Entry struct {
	k {{ .KeyType }}
	v {{ .ValueType }}
}
type {{ .Name }}MapBuilder {{ .Name }}Map

func {{ .Method -}} Make{{ .Name | upper }}Map(ents ...{{ .Name }}Entry) {{ .Name }}Map {
	x := make(map[{{ .KeyType }}]{{ .ValueType }}, len(ents))
	for _, y := range ents {
		x[y.k] = y.v
	}
	return {{ .Name }}Map{x}
}
func {{ .Method -}} Make{{ .Name | upper }}MapEntry(k {{ .KeyType }}, v {{ .ValueType }}) {{ .Name }}Entry {
	return {{ .Name }}Entry{k, v}
}
func {{ .Method -}} Start{{ .Name | upper }}Map(sizeHint int) {{ .Name }}MapBuilder {
	return {{ .Name }}MapBuilder{make(map[{{ .KeyType }}]{{ .ValueType }}, sizeHint)}
}
func (b *{{ .Name }}MapBuilder) Append(k {{ .KeyType }}, v {{ .ValueType }}) {
	b.x[k] = v
}
func (b *{{ .Name }}MapBuilder) Finish() {{ .Name }}Map {
	v := *b
	b.x = nil
	return {{ .Name }}Map(v)
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
	return startCase(nostar(li.ValueType))
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
type {{ .Name }}List struct {
	x []{{ .ValueType }}
}
type {{ .Name }}ListBuilder {{ .Name }}List

func {{ .Method -}} Make{{ .Name | upper }}List(ents ...{{ .ValueType }}) {{ .Name }}List {
	x := make([]{{ .ValueType }}, len(ents))
	copy(x, ents)
	return {{ .Name }}List{x}
}
func {{ .Method -}} Start{{ .Name | upper }}List(sizeHint int) {{ .Name }}ListBuilder {
	return {{ .Name }}ListBuilder{make([]{{ .ValueType }}, 0, sizeHint)}
}
func (b *{{ .Name }}ListBuilder) Append(v {{ .ValueType }}) {
	b.x = append(b.x, v)
}
func (b *{{ .Name }}ListBuilder) Finish() {{ .Name }}List {
	v := *b
	b.x = nil
	return {{ .Name }}List(v)
}
`
