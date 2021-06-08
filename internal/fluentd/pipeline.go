package fluentd

import (
	"text/template"
)

type Pipeline struct {
	InLabel
	SubElements []Element
	Desc        string
}

func (p Pipeline) Name() string {
	return "pipeline"
}

func (p Pipeline) Template() string {
	return `{{define "` + p.Name() + `"  -}}
# {{.Desc}}
<label @{{.InLabel}}>
{{generate .Elements | indent 2}}
</label>
{{end}}
`
}

func (p Pipeline) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(p.Template()))
}

func (p Pipeline) Data() interface{} {
	return p
}
func (p Pipeline) Elements() []Element {
	return p.SubElements
}
