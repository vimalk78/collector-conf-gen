package fluentdforward

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type SharedKey security.SharedKey

func (sk SharedKey) Name() string {
	return "fluentdforwardSharedKeyTemplate"
}

func (sk SharedKey) Template() string {
	return `{{define "` + sk.Name() + `" -}}
<security>
  self_hostname "#{ENV['NODE_NAME']}"
  shared_key "{{.KeyPath}}"
</security>
{{- end}}`
}

func (sk SharedKey) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(sk.Template()))
}

func (sk SharedKey) Data() interface{} {
	return sk
}
