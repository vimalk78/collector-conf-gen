package fluentd

import (
	"text/template"

	"github.com/vimalk78/collector-conf-gen/internal/generator"
)

type Copy struct {
	Labels []string
}

func (c Copy) Name() string {
	return "copySourceTypeToPipeline"
}

func (c Copy) Template() string {
	return `{{define "` + c.Name() + `"  -}}
<match **>
  @type copy
  {{- range $index, $label := .Labels }}
  <store>
    @type relabel
    @label {{ $label }}
  </store>
  {{- end }}
</match>
{{- end}}`
}

func (c Copy) Create(t *template.Template, ct generator.CollectorConfType) *template.Template {
	return template.Must(t.Parse(generator.SelectTemplate(c, ct)))
}

func (c Copy) Data() interface{} {
	return c
}
