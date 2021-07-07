package assembler

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type ContainerLogs struct {
	OutLabel
	Desc         string
	Paths        string
	ExcludePaths string
	PosFile      string
}

func (cl ContainerLogs) Name() string {
	return "inputContainerSourceTemplate"
}

func (cl ContainerLogs) Template() string {
	return `{{define "` + cl.Name() + `" -}}
# {{.Desc}}
<source>
  @type tail
  @id container-input
  path {{.Paths}}
  exclude_path {{.ExcludePaths}}
  pos_file "{{.PosFile}}"
  refresh_interval 5
  rotate_wait 5
  tag kubernetes.*
  read_from_head "true"
  @label @{{.OutLabel}}
  <parse>
    @type multi_format
    <pattern>
      format json
      time_format '%Y-%m-%dT%H:%M:%S.%N%Z'
      keep_time_key true
    </pattern>
    <pattern>
      format regexp
      expression /^(?<time>[^\s]+) (?<stream>stdout|stderr)( (?<logtag>.))? (?<log>.*)$/
      time_format '%Y-%m-%dT%H:%M:%S.%N%:z'
      keep_time_key true
    </pattern>
  </parse>
</source>
{{end}}`
}

func (cl ContainerLogs) TemplateVector() string {
	return `{{define "` + cl.Name() + `" -}}
[sources.{{.OutLabel}}]
  type = "file" # required
  ignore_older_secs = 600
  include = [{{.Paths}}]
  exclude = {{.ExcludePaths}}
  read_from = "beginning"
  data_dir = {{.PosFile}}
{{end}}
`
}

func (cl ContainerLogs) Create(t *template.Template, ct CollectorConfType) *template.Template {
	var tmpl string
	if ct == CollectorConfVector {
		tmpl = cl.TemplateVector()
	} else {
		tmpl = cl.Template()
	}
	return template.Must(t.Parse(tmpl))
}

func (cl ContainerLogs) Data() interface{} {
	return cl
}
