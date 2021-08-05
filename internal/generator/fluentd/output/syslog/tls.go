package syslog

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type TLS security.TLS

func (t TLS) Name() string {
	return "syslogTLSTemplate"
}

func (t TLS) Template() string {
	if t {
		return `{{define "syslogTLSTemplate" -}}
tls true
{{end}}`
	}
	return `{{define "syslogTLSTemplate" -}}
tls false
{{end}}`
}

func (t TLS) Create(tp *template.Template) *template.Template {
	return template.Must(tp.Parse(t.Template()))
}

func (t TLS) Data() interface{} {
	return t
}
