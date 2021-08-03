package cloudwatch

import "text/template"

type AWSKey struct {
	KeyIDPath string
	KeyPath   string
}

func (a AWSKey) Name() string {
	return "awsKeyTemplate"
}

func (a AWSKey) Template() string {
	return `{{define "` + a.Name() + `" -}}
aws_key_id "#{open('{{ .KeyIDPath }}','r') do |f|f.read end}"
aws_sec_key "#{open('{{ .KeyPath }}','r') do |f|f.read end}"
{{end}}`
}

func (a AWSKey) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(a.Template()))
}

func (a AWSKey) Data() interface{} {
	return a
}
