package fluentd

import "text/template"

type Relabel struct {
	Desc string
	OutLabel
	Pattern string
}

func (r Relabel) Name() string {
	return "relabel"
}

func (r Relabel) Template() string {
	return `{{define "` + r.Name() + `"  -}}
# {{.Desc}}
<match {{.Pattern}}>
  @type relabel
  @label @{{.OutLabel}}
</match>
{{end}}
`
}

func (r Relabel) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(r.Template()))
}

func (r Relabel) Data() interface{} {
	return r
}
