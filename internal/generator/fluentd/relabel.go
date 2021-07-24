package fluentd

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Relabel struct {
	OutLabel
}

func (r Relabel) Name() string {
	return "relabel"
}

func (r Relabel) Template() string {
	return `{{define "` + r.Name() + `"  -}}
@type relabel
@label {{.OutLabel}}
{{- end}}`
}

func (r Relabel) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(r.Template()))
}

func (r Relabel) Data() interface{} {
	return r
}
