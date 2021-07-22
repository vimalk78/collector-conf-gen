package kafka

import (
	"text/template"
)

type SaslOverSSL bool

func (s SaslOverSSL) Name() string {
	return "kafkaSaslOverSSLTemplate"
}

func (s SaslOverSSL) Template() string {
	enabled := `{{define "kafkaSaslOverSSLTemplate" -}}
sasl_over_ssl true
{{- end}}
`
	disabled := `{{define "kafkaSaslOverSSLTemplate" -}}
sasl_over_ssl false
{{- end}}
`
	if s {
		return enabled
	}
	return disabled
}

func (s SaslOverSSL) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(s.Template()))
}

func (s SaslOverSSL) Data() interface{} {
	return s
}
