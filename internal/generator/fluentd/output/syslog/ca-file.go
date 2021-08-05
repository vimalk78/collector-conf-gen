package syslog

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type CAFile security.CAFile

func (ca CAFile) Name() string {
	return "syslogCAFileTemplate"
}

func (ca CAFile) Template() string {
	return `{{define "` + ca.Name() + `" -}}
ca_file {{.CAFilePath}}
{{- end}}
`
}

func (ca CAFile) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(ca.Template()))
}

func (ca CAFile) Data() interface{} {
	return ca
}
