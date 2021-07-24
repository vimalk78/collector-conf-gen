package fluentdforward

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type TLS security.TLS

func (t TLS) Name() string {
	return "fluentdforwardTLSTemplate"
}

func (t TLS) Template() string {
	https := `{{define "fluentdforwardTLSTemplate" -}}
transport tls
tls_verify_hostname false
tls_version 'TLSv1_2'
{{- end}}
`
	http := `{{define "fluentdforwardTLSTemplate" -}}
tls_insecure_mode true
{{- end}}
`
	if t {
		return https
	}
	return http
}

func (t TLS) Create(tp *template.Template) *template.Template {
	return template.Must(tp.Parse(t.Template()))
}

func (t TLS) Data() interface{} {
	return t
}
