package fluentd

import "text/template"

func NewStdout() Element {
	return Stdout{
		Pattern: "*",
	}
}

// --- output Stdout
type Stdout struct {
	InLabel
	Pattern string
	Desc    string
}

func (os Stdout) Name() string {
	return "outputStdout"
}

func (os Stdout) Template() string {
	return `{{define "` + os.Name() + `"  -}}
# {{.Desc}}
<match {{if .Pattern}}{{.Pattern}}{{else}}*{{end}}>
  @type stdout
</match>
{{end}}`
}

func (os Stdout) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(os.Template()))
}

func (os Stdout) Data() interface{} {
	return os
}
