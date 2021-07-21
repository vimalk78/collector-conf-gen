package fluentd

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Match struct {
	MatchTags string
	Elements  []Element
}

func (m Match) Name() string {
	return "matchTemplate"
}

func (m Match) Template() string {
	return `{{define "` + m.Name() + `"  -}}
<match {{.MatchTags}}>
{{compose .Elements | indent 2}}
</match>
{{end}}`
}

func (m Match) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(m.Template()))
}

func (m Match) Data() interface{} {
	return m
}
