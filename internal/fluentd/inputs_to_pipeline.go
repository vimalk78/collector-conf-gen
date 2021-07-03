package fluentd

import (
	"fmt"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
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

func (g *Generator) SourceTypeToPipeline(sourceType string, spec *logging.ClusterLogForwarderSpec) Element {
	c := Copy{}
	for _, pipeline := range spec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if inRef == sourceType {
				c.Labels = append(c.Labels, pipeline.Name)
			}
		}
	}
	if len(c.Labels) == 0 {
		return _Nil
	}
	if len(c.Labels) == 1 {
		return FromLabel{
			Desc:    fmt.Sprintf("Sending %s source type to pipeline", sourceType),
			InLabel: sourceTypeLabelName(sourceType),
			SubElements: []Element{
				Relabel{
					MatchTags: "**",
					OutLabel:  labelName(c.Labels[0]),
				},
			},
		}
	}
	return FromLabel{
		Desc:    fmt.Sprintf("Copying %s source type to pipeline", sourceType),
		InLabel: sourceTypeLabelName(sourceType),
		SubElements: []Element{
			c,
		},
	}
}

func (g *Generator) InputsToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return MergeElements(
		g.ApplicationToPipeline(spec),
		g.InfraToPipeline(spec),
		g.AuditToPipeline(spec),
	)
}

func (g *Generator) ApplicationToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	userDefined := spec.InputMap()
	p := ApplicationsToPipelines{}
	c := Copy{}
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
								Pipeline: pipeline.Name,
							}
						}
						a.Namespaces = app.Namespaces
					}
					if app.Selector != nil && len(app.Selector.MatchLabels) != 0 {
						if a == nil {
							a = &ApplicationToPipeline{
								Pipeline: pipeline.Name,
							}
						}
						a.Labels = LabelsKV(app.Selector)
					}
					if a != nil {
						p = append(p, *a)
					} else {
						c.Labels = append(c.Labels, pipeline.Name)
					}
				}
			} else if inRef == logging.InputNameApplication {
				c.Labels = append(c.Labels, pipeline.Name)
			}
		}
	}
	if len(p) == 0 {
		return []Element{
			g.SourceTypeToPipeline(logging.InputNameApplication, spec),
		}
	}
	l := FromLabel{
		Desc:    "Copying unrouted \"application\" to pipelines",
		InLabel: sourceTypeLabelName("APPLICATION_ALL"),
		SubElements: []Element{
			c,
		},
	}
	if len(c.Labels) != 0 {
		p = append(p, ApplicationToPipeline{
			Pipeline: "_APPLICATION_ALL",
		})
		return []Element{
			p,
			l,
		}
	} else {
		return []Element{
			p,
		}
	}
}

func (g *Generator) AuditToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return []Element{
		g.SourceTypeToPipeline(logging.InputNameAudit, spec),
	}
}

func (g *Generator) InfraToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return []Element{
		g.SourceTypeToPipeline(logging.InputNameInfrastructure, spec),
	}
}
