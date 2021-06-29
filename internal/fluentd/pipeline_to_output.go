package fluentd

import (
	"encoding/json"
	"fmt"
	"sort"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

type PipelineToOutputs struct {
	Desc      string
	Pipeline  string
	Outputs   []string
	Labels    []Element
	JsonParse []Element
}

func (p PipelineToOutputs) Name() string {
	return "pipelineToPutput"
}

func (p PipelineToOutputs) Template() string {
	return `{{define "` + p.Name() + `"  -}}
<label {{labelName .Pipeline}}>
{{- with $x := generate .Labels}}
{{$x | indent 2 -}}
{{- end}}
{{- with $x := generate .JsonParse}}
{{$x |indent 2 -}}
{{- end}}
  <match **>
    @type copy
    {{- range $index, $output := .Outputs }}
    <store>
      @type relabel
      @label {{labelName $output }}
    </store>
    {{- end }}
  </match>
</label>
{{- end}}`
}

func (p PipelineToOutputs) Data() interface{} {
	return p
}

func (p PipelineToOutputs) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(p.Template()))
}

var PipelineLabels = `
{{define "PipelineLabels" -}}
# {{.Desc}}
<filter **>
  @type record_transformer
  <record>
    openshift { "labels": %s }
  </record>
</filter>
{{- end}}`

var JsonParseTemplate = `{{define "JsonParse" -}}
# {{.Desc}}
<filter **>
  @type parser
  key_name message
  reserve_data yes
  hash_value_field structured
  <parse>
    @type json
    json_parser oj
  </parse>
</filter>
{{- end}}`

func (g *Generator) PipelineToOutputs(spec *logging.ClusterLogForwarderSpec) []Element {
	var e []Element = []Element{}
	pipelines := spec.Pipelines
	sort.Slice(pipelines, func(i, j int) bool {
		return pipelines[i].Name < pipelines[j].Name
	})
	for _, p := range pipelines {
		po := PipelineToOutputs{
			Pipeline:  p.Name,
			Outputs:   p.OutputRefs,
			JsonParse: _Nils,
			Labels:    _Nils,
		}
		if p.Labels != nil && len(p.Labels) != 0 {
			// ignoring error, because pre-check stage already checked if Labels can be marshalled
			s, _ := json.Marshal(p.Labels)
			po.Labels = []Element{
				ConfLiteral{
					Desc:         "Add User Defined labels to the output record",
					TemplateName: "PipelineLabels",
					TemplateStr:  fmt.Sprintf(PipelineLabels, string(s)),
				},
			}
		}
		if p.Parse == "json" {
			po.JsonParse = []Element{
				ConfLiteral{
					Desc:         "Parse the logs into json",
					TemplateName: "JsonParse",
					TemplateStr:  JsonParseTemplate,
				},
			}
		}
		e = append(e, po)
	}
	return e
}
