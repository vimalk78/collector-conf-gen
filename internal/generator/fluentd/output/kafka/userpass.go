package kafka

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type UserNamePass security.UserNamePass

func (up UserNamePass) Name() string {
	return "kafkaUsernamePasswordTemplate"
}

func (up UserNamePass) Template() string {
	return `{{define "` + up.Name() + `" -}}
sasl_plain_username "#{File.exists?('{{.UsernamePath}}') ? open('{{.UsernamePath}}','r') do |f|f.read end : ''}"
sasl_plain_password "#{File.exists?('{{.PasswordPath}}') ? open('{{.PasswordPath}}','r') do |f|f.read end : ''}"
{{- end}}
`
}

func (up UserNamePass) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(up.Template()))
}

func (up UserNamePass) Data() interface{} {
	return up
}
