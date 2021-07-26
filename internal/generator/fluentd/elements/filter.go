package elements

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Filter struct {
	Desc      string
	MatchTags string
	Element   Element
}

func (f Filter) Name() string {
	return "filterTemplate"
}

func (f Filter) Template() string {
	return `{{define "` + f.Name() + `" -}}
{{- if .Desc}}
#{{.Desc}}
{{- end}}
<filter {{.MatchTags}}>
{{compose_one .Element | indent 2}}
</filter>
{{- end}}
`
}

func (f Filter) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(f.Template()))
}

func (f Filter) Data() interface{} {
	return f
}
