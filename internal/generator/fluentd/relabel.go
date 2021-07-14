package fluentd

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Relabel struct {
	Desc string
	OutLabel
	MatchTags string
}

func (r Relabel) Name() string {
	return "relabel"
}

func (r Relabel) Template() string {
	return `{{define "` + r.Name() + `"  -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
<match {{.MatchTags}}>
  @type relabel
  @label {{.OutLabel}}
</match>
{{- end}}`
}

func (r Relabel) Create(t *template.Template, ct CollectorConfType) *template.Template {
	return template.Must(t.Parse(SelectTemplate(r, ct)))
}

func (r Relabel) Data() interface{} {
	return r
}
