package elements

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Match struct {
	Desc         string
	MatchTags    string
	MatchElement Element
}

func (m Match) Name() string {
	return "matchTemplate"
}

func (m Match) Template() string {
	return `{{define "` + m.Name() + `"  -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
<match {{.MatchTags}}>
{{compose_one .MatchElement | indent 2}}
</match>
{{- end}}`
}

func (m Match) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(m.Template()))
}

func (m Match) Data() interface{} {
	return m
}
