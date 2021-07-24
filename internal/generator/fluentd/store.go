package fluentd

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Store struct {
	Element Element
}

func (s Store) Name() string {
	return "storeTemplate"
}

func (s Store) Template() string {
	return `{{define "` + s.Name() + `" -}}
<store>
{{compose_one .Element| indent 2}}
</store>
{{- end}}
`
}

func (s Store) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(s.Template()))
}

func (s Store) Data() interface{} {
	return s
}
