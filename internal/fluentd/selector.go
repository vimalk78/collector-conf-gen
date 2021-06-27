package fluentd

import "text/template"

type CopySourceTypeToPipeline struct {
	SourceType string
	Pipelines  []string
	Desc       string
}

func (c CopySourceTypeToPipeline) Name() string {
	return "copySourceTypeToPipeline"
}

func (c CopySourceTypeToPipeline) Template() string {
	return `{{define "` + c.Name() + `"  -}}
# {{.Desc}}
<label {{sourceTypelabelName .SourceType}}>
  <match **>
    @type copy
    {{- range $index, $pipeline := .Pipelines }}
    <store>
      @type relabel
      @label {{labelName $pipeline }}
    </store>
    {{- end }}
  </match>
</label>
{{end}}`
}

func (c CopySourceTypeToPipeline) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(c.Template()))
}

func (c CopySourceTypeToPipeline) Data() interface{} {
	return c
}

type ApplicationToPipeline struct {
	// Labels is an array of "<key>:<value>" strings
	Labels     []string
	Namespaces []string
	Pipeline   string
}

type ApplicationsToPipelines []ApplicationToPipeline

func (a ApplicationsToPipelines) Name() string {
	return "applicationToPipeline"
}

func (a ApplicationsToPipelines) Template() string {
	return `{{define "` + a.Name() + `"  -}}
# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
	{{- range $index, $a := .}}
    <route>
      @label {{labelName $a.Pipeline}}
      <match>
        {{- if $a.Namespaces}}
        namespaces {{comma_separated $a.Namespaces}}
		{{- end}}
        {{- if $a.Labels}}
        labels {{comma_separated $a.Labels }}
		{{- end}}
      </match>
    </route>
	{{- end}}
  </match>
</label>
{{end}}`
}

func (a ApplicationsToPipelines) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(a.Template()))
}

func (a ApplicationsToPipelines) Data() interface{} {
	return a
}
