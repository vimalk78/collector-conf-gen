package fluentd

import (
	"encoding/json"
	"fmt"
	"sort"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type PipelineToOutputs_ struct {
	Desc      string
	Pipeline  string
	Labels    []Element
	JsonParse []Element
	ToOutputs []Element
}

func (p PipelineToOutputs_) Name() string {
	return "pipelineToPutput"
}

func (p PipelineToOutputs_) Template() string {
	return `{{define "` + p.Name() + `"  -}}
# {{.Desc}}
<label {{.Pipeline}}>
{{- with $x := compose .Labels}}
{{$x |indent 2 -}}
{{- end}}
{{- with $x := compose .JsonParse}}
{{$x |indent 2 -}}
{{- end}}
{{compose .ToOutputs| indent 2}}
</label>
{{- end}}`
}

func (p PipelineToOutputs_) Data() interface{} {
	return p
}

func (p PipelineToOutputs_) Create(t *template.Template) *template.Template {
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

func PipelineToOutputs(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	var e []Element = []Element{}
	pipelines := spec.Pipelines
	sort.Slice(pipelines, func(i, j int) bool {
		return pipelines[i].Name < pipelines[j].Name
	})
	for _, p := range pipelines {
		po := PipelineToOutputs_{
			Desc:      fmt.Sprintf("Copying pipeline %s to outputs", p.Name),
			Pipeline:  labelName(p.Name),
			JsonParse: Nils,
			Labels:    Nils,
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
		switch len(p.OutputRefs) {
		case 0:
			// should not happen
		case 1:
			po.ToOutputs = []Element{
				Match{
					MatchTags: "**",
					Elements: []Element{
						Relabel{
							OutLabel: labelName(p.OutputRefs[0]),
						},
					},
				},
			}
		default:
			po.ToOutputs = []Element{
				Copy{
					Labels: labelNames(p.OutputRefs),
				},
			}
		}
		e = append(e, po)
	}
	return e
}
