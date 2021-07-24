package fluentdforward

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type CAFile security.CAFile

func (ca CAFile) Name() string {
	return "fluentdforwardCAFileTemplate"
}

func (ca CAFile) Template() string {
	return `{{define "` + ca.Name() + `" -}}
tls_cert_path {{.CAFilePath}}
{{- end}}
`
}

func (ca CAFile) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(ca.Template()))
}

func (ca CAFile) Data() interface{} {
	return ca
}
