package elements

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
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
{{if .Desc -}}
# {{.Desc}}
{{end -}}
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
