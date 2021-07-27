package elasticsearch

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
)

type UserNamePass security.UserNamePass

func (up UserNamePass) Name() string {
	return "elasticsearchUsernamePasswordTemplate"
}

func (up UserNamePass) Template() string {
	return `{{define "` + up.Name() + `" -}}
username "#{File.exists?({{.UsernamePath}}) ? open({{.UsernamePath}},'r') do |f|f.read end : ''}"
password "#{File.exists?({{.PasswordPath}}) ? open({{.PasswordPath}},'r') do |f|f.read end : ''}"
{{- end}}
`
}

func (up UserNamePass) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(up.Template()))
}

func (up UserNamePass) Data() interface{} {
	return up
}
