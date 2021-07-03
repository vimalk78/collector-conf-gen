package fluentd

import "text/template"

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
    @label {{labelName $label }}
  </store>
  {{- end }}
</match>
{{- end}}`
}

func (c Copy) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(c.Template()))
}

func (c Copy) Data() interface{} {
	return c
}
