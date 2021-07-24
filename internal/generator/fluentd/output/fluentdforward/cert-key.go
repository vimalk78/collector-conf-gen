package fluentdforward

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type TLSKeyCert security.TLSKeyCert

func (kc TLSKeyCert) Name() string {
	return "fluentdforwardCertKeyTemplate"
}

func (kc TLSKeyCert) Template() string {
	return `{{define "` + kc.Name() + `" -}}
tls_client_private_key_path {{.KeyPath}}
tls_client_cert_path {{.CertPath}}
{{- end}}`
}

func (kc TLSKeyCert) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(kc.Template()))
}

func (kc TLSKeyCert) Data() interface{} {
	return kc
}
