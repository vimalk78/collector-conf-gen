package fluentd

import (
	"fmt"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
)

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
      @label {{$a.Pipeline}}
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
{{- end}}`
}

func (a ApplicationsToPipelines) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(a.Template()))
}

func (a ApplicationsToPipelines) Data() interface{} {
	return a
}

func SourceTypeToPipeline(sourceType string, spec *logging.ClusterLogForwarderSpec, op *Options) Element {
	srcTypePipeline := []string{}
	for _, pipeline := range spec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if inRef == sourceType {
				srcTypePipeline = append(srcTypePipeline, pipeline.Name)
			}
		}
	}
	if Clo.IncludeLegacyForwardConfig(*op) {
		srcTypePipeline = append(srcTypePipeline, LegacySecureforward)
	}
	if Clo.IncludeLegacySyslogConfig(*op) {
		srcTypePipeline = append(srcTypePipeline, LegacySyslog)
	}
	switch len(srcTypePipeline) {
	case 0:
		return Nil
	case 1:
		return FromLabel{
			Desc:    fmt.Sprintf("Sending %s source type to pipeline", sourceType),
			InLabel: helpers.SourceTypeLabelName(sourceType),
			SubElements: []Element{
				Match{
					MatchTags: "**",
					MatchElement: Relabel{
						OutLabel: helpers.LabelName(srcTypePipeline[0]),
					},
				},
			},
		}
	default:
		return FromLabel{
			Desc:    fmt.Sprintf("Copying %s source type to pipeline", sourceType),
			InLabel: helpers.SourceTypeLabelName(sourceType),
			SubElements: []Element{
				Match{
					MatchTags: "**",
					MatchElement: Copy{
						Stores: CopyToLabels(helpers.LabelNames(srcTypePipeline)),
					},
				},
			},
		}
	}
}

func InputsToPipeline(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return MergeElements(
		AppToPipeline(spec, o),
		InfraToPipeline(spec, o),
		AuditToPipeline(spec, o),
	)
}

func AppToPipeline(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	userDefined := spec.InputMap()
	// routed by namespace, or labels
	routedPipelines := ApplicationsToPipelines{}
	unRoutedPipelines := []string{}
	for _, pipeline := range spec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if input, ok := userDefined[inRef]; ok {
				// user defined input
				if input.Application != nil {
					app := input.Application
					var a *ApplicationToPipeline = nil
					if len(app.Namespaces) != 0 {
						if a == nil {
							a = &ApplicationToPipeline{
								Pipeline: helpers.LabelName(pipeline.Name),
							}
						}
						a.Namespaces = app.Namespaces
					}
					if app.Selector != nil && len(app.Selector.MatchLabels) != 0 {
						if a == nil {
							a = &ApplicationToPipeline{
								Pipeline: helpers.LabelName(pipeline.Name),
							}
						}
						a.Labels = helpers.LabelsKV(app.Selector)
					}
					if a != nil {
						routedPipelines = append(routedPipelines, *a)
					} else {
						unRoutedPipelines = append(unRoutedPipelines, pipeline.Name)
					}
				}
			} else if inRef == logging.InputNameApplication {
				unRoutedPipelines = append(unRoutedPipelines, pipeline.Name)
			}
		}
	}
	if len(routedPipelines) == 0 {
		return []Element{
			SourceTypeToPipeline(logging.InputNameApplication, spec, o),
		}
	}
	fmt.Printf("unRoutedPipelines: %v\n", unRoutedPipelines)
	switch len(unRoutedPipelines) {
	case 0:
		return []Element{
			routedPipelines,
		}
	case 1:
		routedPipelines = append(routedPipelines, ApplicationToPipeline{
			Pipeline: helpers.SourceTypeLabelName("APPLICATION_ALL"),
		})
		return []Element{
			routedPipelines,
			FromLabel{
				Desc:    "Sending unrouted application to pipelines",
				InLabel: helpers.SourceTypeLabelName("APPLICATION_ALL"),
				SubElements: []Element{
					Match{
						MatchTags: "**",
						MatchElement: Relabel{
							OutLabel: helpers.LabelName(unRoutedPipelines[0]),
						},
					},
				},
			},
		}
	default:
		routedPipelines = append(routedPipelines, ApplicationToPipeline{
			Pipeline: helpers.SourceTypeLabelName("APPLICATION_ALL"),
		})
		return []Element{
			routedPipelines,
			FromLabel{
				Desc:    "Copying unrouted application to pipelines",
				InLabel: helpers.SourceTypeLabelName("APPLICATION_ALL"),
				SubElements: []Element{
					Match{
						MatchTags: "**",
						MatchElement: Copy{
							Stores: CopyToLabels(helpers.LabelNames(unRoutedPipelines)),
						},
					},
				},
			},
		}
	}
}

func AuditToPipeline(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return []Element{
		SourceTypeToPipeline(logging.InputNameAudit, spec, o),
	}
}

func InfraToPipeline(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return []Element{
		SourceTypeToPipeline(logging.InputNameInfrastructure, spec, o),
	}
}
