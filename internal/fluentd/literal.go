package fluentd

import "text/template"

type ConfLiteral struct {
	TemplateName string
	Desc         string
	InLabel
	OutLabel
	Pattern     string
	TemplateStr string
}

func (b ConfLiteral) Name() string {
	return b.TemplateName
}

func (b ConfLiteral) Template() string {
	return b.TemplateStr
}

func (b ConfLiteral) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(b.Template()))
}

func (b ConfLiteral) Data() interface{} {
	return b
}
