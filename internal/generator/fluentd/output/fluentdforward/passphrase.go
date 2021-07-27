package fluentdforward

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type Passphrase security.Passphrase

func (p Passphrase) Name() string {
	return "passphraseTemplate"
}

func (p Passphrase) Template() string {
	return `{{define "` + p.Name() + `" -}}
tls_client_private_key_passphrase "#{File.exists?({{.PassphrasePath}}) ? open({{.PassphrasePath}},'r') do |f|f.read end : ''}" 
{{- end}}
`
}

func (p Passphrase) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(p.Template()))
}

func (p Passphrase) Data() interface{} {
	return p
}
