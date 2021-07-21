package elasticsearch

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type TLSKeyCert security.TLSKeyCert

func (kc TLSKeyCert) Name() string {
	return "elasticsearchCertKeyTemplate"
}

func (kc TLSKeyCert) Template() string {
	return `{{define "` + kc.Name() + `" -}}
client_key {{.KeyPath}}
client_cert {{.CertPath}}
{{- end}}`
}

func (kc TLSKeyCert) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(kc.Template()))
}

func (kc TLSKeyCert) Data() interface{} {
	return kc
}
