package elements

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Copy struct {
	Stores []Element
}

func (c Copy) Name() string {
	return "copySourceTypeToPipeline"
}

func (c Copy) Template() string {
	return `{{define "` + c.Name() + `"  -}}
@type copy
{{compose .Stores}}
{{- end}}`
}

func (c Copy) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(c.Template()))
}

func (c Copy) Data() interface{} {
	return c
}

func CopyToLabels(labels []string) []Element {
	s := []Element{}
	for _, l := range labels {
		s = append(s, Store{
			Element: Relabel{
				OutLabel: l,
			},
		})
	}
	return s
}
