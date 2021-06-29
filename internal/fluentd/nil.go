package fluentd

import "text/template"

type NilElement int

func (r NilElement) Name() string {
	return "nilElement"
}

func (r NilElement) Template() string {
	return `{{- define "` + r.Name() + `"  -}}{{- end -}}`
}

func (n NilElement) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(n.Template()))
}

func (n NilElement) Data() interface{} {
	return n
}

var _Nil NilElement
var _Nils []Element = []Element{_Nil}
