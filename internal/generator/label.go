package generator

import (
	"text/template"
)

type FromLabel struct {
	InLabel
	SubElements []Element
	Desc        string
}

func (f FromLabel) Name() string {
	return "pipeline"
}

func (f FromLabel) Template() string {
	return `{{define "` + f.Name() + `"  -}}
# {{.Desc}}
<label {{.InLabel}}>
{{compose .Elements | indent 2}}
</label>
{{end}}`
}

func (f FromLabel) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(f.Template()))
}

func (f FromLabel) Data() interface{} {
	return f
}
func (f FromLabel) Elements() []Element {
	return f.SubElements
}

type Pipeline = FromLabel
