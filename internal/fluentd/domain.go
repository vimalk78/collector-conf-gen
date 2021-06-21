package fluentd

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

var clo CLO

type Generator struct {
	// spec
	// routemap
}

func MakeGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) MakeLoggingConf(spec *logging.ClusterLogForwarderSpec) []Section {
	return []Section{
		{
			g.Sources(spec),
			"Set of all input sources",
		},
		{
			g.Metrics(),
			"Section to add measurement, and dispatch to Concat or Ingress pipelines",
		},
		{
			g.Concat(),
			"Concat pipeline",
		},
		{
			g.Ingress(spec),
			"Ingress pipeline",
		},
		{
			g.InputsToPipeline(spec),
			"Inputs go to pipelines",
		},
	}
}

func (g *Generator) Sources(spec *logging.ClusterLogForwarderSpec) []Element {
	var el []Element = make([]Element, 0)
	el = append(el, PrometheusMonitor{
		Desc: "Prometheus Monitoring",
	})
	types := clo.GatherSources(spec)
	if types.Has(logging.InputNameApplication) {
		el = append(el,
			ContainerLogs{
				OutLabel: "MEASURE",
				Desc:     "Logs from containers (including openshift containers)",
			})
	}
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el,
			JournalLog{
				OutLabel: "MEASURE",
				Desc:     "Logs from linux journal",
			})
	}
	if types.Has(logging.InputNameAudit) {
		el = append(el,
			HostAuditLog{
				OutLabel: "MEASURE",
				Desc:     "Logs from host audit",
			},
			K8sAuditLog{
				OutLabel: "MEASURE",
				Desc:     "Logs from kubernetes audit",
			},
			OpenshiftAuditLog{
				OutLabel: "MEASURE",
				Desc:     "Logs from openshift audit",
			})
	}
	return el
}

func (g *Generator) Metrics() []Element {
	return []Element{
		Pipeline{
			InLabel: "MEASURE",
			Desc:    "Increment Prometheus metrics",
			SubElements: []Element{
				ConfLiteral{
					TemplateName: "EmitMetrics",
					Desc:         "xx",
					TemplateStr:  EmitMetrics,
				},
				Relabel{
					Desc:     "Journal Logs go to INGRESS pipeline",
					Pattern:  "journal",
					OutLabel: "INGRESS",
				},
				Relabel{
					Desc:     "Audit Logs go to INGRESS pipeline",
					Pattern:  "*.audit.log",
					OutLabel: "INGRESS",
				},
				Relabel{
					Desc:     "Kubernetes Logs go to CONCAT pipeline",
					Pattern:  "kubernetes.**",
					OutLabel: "CONCAT",
				},
			},
		},
	}
}

func (g *Generator) Concat() []Element {
	return []Element{
		Pipeline{
			InLabel: "CONCAT",
			Desc:    "Concat log lines of container logs",
			SubElements: []Element{
				ConfLiteral{
					Desc:         "Concat container lines",
					TemplateName: "concatLines",
					TemplateStr:  ConcatLines,
				},
				Relabel{
					Desc:     "Kubernetes Logs go to INGRESS pipeline",
					Pattern:  "kubernetes.**",
					OutLabel: "INGRESS",
				},
			},
		},
	}
}

func (g *Generator) Ingress(spec *logging.ClusterLogForwarderSpec) []Element {
	return []Element{
		Pipeline{
			InLabel: "INGRESS",
			Desc:    "Concat log lines of container logs",
			SubElements: MergeElements([]Element{
				ConfLiteral{
					Desc:         "Set Encodeing",
					TemplateName: "setEncoding",
					TemplateStr:  SetEncoding,
				},
				ConfLiteral{
					Desc:         "Filter out PRIORITY from journal logs",
					TemplateName: "filterJournalPRIORITY",
					TemplateStr:  FilterJournalPRIORITY,
				},
				ConfLiteral{
					Desc:         "Retag Journal logs to specific tags",
					OutLabel:     "INGRESS",
					TemplateName: "retagJournal",
					TemplateStr:  RetagJournalLogs,
				},
				ConfLiteral{
					Desc:         "Invoke kubernetes apiserver to get kunbernetes metadata",
					TemplateName: "kubernetesMetadata",
					TemplateStr:  KubernetesMetadataPlugin,
				},
				ConfLiteral{
					Desc:         "Parse Json fields for container, journal and eventrouter logs",
					TemplateName: "parseJsonFields",
					TemplateStr:  ParseJsonFields,
				},
				ConfLiteral{
					Desc:         "Clean kibana log fields",
					TemplateName: "cleanKibanaLogs",
					TemplateStr:  CleanKibanaLogs,
				},
				ConfLiteral{
					Desc:         "Fix level field in audit logs",
					TemplateName: "fixAuditLevel",
					TemplateStr:  FixAuditLevel,
				},
				ConfLiteral{
					Desc:         "Viaq Data Model: The big bad viaq model.",
					TemplateName: "viaqDataModel",
					TemplateStr:  ViaQDataModel,
				},
				ConfLiteral{
					Desc:         "Generate elasticsearch id",
					TemplateName: "genElasticsearchID",
					TemplateStr:  GenElasticsearchID,
				},
			},
				g.SelectLogTypePipeline(spec)),
		},
	}
}

func (g *Generator) SelectLogTypePipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	var el []Element = make([]Element, 0)
	types := clo.GatherSources(spec)
	ApplicationPattern := "kubernetes.**"
	if types.Has(logging.InputNameApplication) {
		el = append(el, Relabel{
			Desc:     "Dont discard Application logs",
			Pattern:  ApplicationPattern,
			OutLabel: "_APPLICATION",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Application logs",
			Pattern:      ApplicationPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	InfraPattern := "**_default_** **_kube-*_** **_openshift-*_** **_openshift_** journal.** system.var.log**"
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el, Relabel{
			Desc:     "Dont discard Infrastructure logs",
			Pattern:  InfraPattern,
			OutLabel: "_INFRASTRUCTURE",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Infrastructure logs",
			Pattern:      InfraPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	AuditPattern := "linux-audit.log** k8s-audit.log** openshift-audit.log**"
	if types.Has(logging.InputNameAudit) {
		el = append(el, Relabel{
			Desc:     "Dont discard Audit logs",
			Pattern:  AuditPattern,
			OutLabel: "_AUDIT",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Audit logs",
			Pattern:      AuditPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	el = append(el, ConfLiteral{
		Desc:         "Send any remaining unmatched tags to stdout",
		TemplateName: "toStdout",
		TemplateStr: `
{{define "toStdout" -}}
# {{.Desc}}
<match **>
 @type stdout
</match>
{{end -}}`,
	})
	return el
}

func (g *Generator) InputsToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return MergeElements(
		g.ApplicationToPipeline(spec),
		g.InfraToPipeline(spec),
		g.AuditToPipeline(spec),
	)
}

func (g *Generator) AuditToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return []Element{
		g.CopySourceTypeToPipeline(logging.InputNameAudit, spec),
	}
}

func (g *Generator) InfraToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	return []Element{
		g.CopySourceTypeToPipeline(logging.InputNameInfrastructure, spec),
	}
}

func (g *Generator) CopySourceTypeToPipeline(sourceType string, spec *logging.ClusterLogForwarderSpec) Element {
	c := CopySourceTypeToPipeline{
		SourceType: sourceType,
		Desc:       fmt.Sprintf("Copying %s source type to pipeline", sourceType),
	}
	for _, pipeline := range spec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if inRef == sourceType {
				c.Pipelines = append(c.Pipelines, pipeline.Name)
			}
		}
	}
	if len(c.Pipelines) == 0 {
		return _Nil
	}
	return c
}

func (g *Generator) ApplicationToPipeline(spec *logging.ClusterLogForwarderSpec) []Element {
	//s := InputSelector{}
	needLabelRouter := false
	userDefined := spec.InputMap()
	for _, ud := range userDefined {
		if ud.Application != nil && (ud.Application.Selector != nil || len(ud.Application.Namespaces) != 0) {
			needLabelRouter = true
		}
	}
	if !needLabelRouter {
		return []Element{
			g.CopySourceTypeToPipeline(logging.InputNameApplication, spec),
		}
	}
	p := ApplicationsToPipelines{}
	for _, pipeline := range spec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if input, ok := userDefined[inRef]; ok {
				// user defined input
				if input.Application != nil && (input.Application.Selector != nil || len(input.Application.Namespaces) != 0) {
					p = append(p, ApplicationToPipeline{
						Pipeline:   pipeline.Name,
						Namespaces: input.Application.Namespaces,
						Labels:     LabelsKV(input.Application.Selector),
					})
				}
				if input.Application != nil {
					fmt.Println("user defined app", pipeline.Name)
				} else {
					// no Namespace or Labels, consider as "application"
					fmt.Println("user defined other", pipeline.Name)
				}
			} else {
				fmt.Println("using default", pipeline.Name)
			}
		}
	}
	return []Element{
		p,
	}
}

/**
 These unexported methods are copied frm CLO
**/
type CLO int

//GatherSources collects the set of unique source types and namespaces
func (CLO) GatherSources(forwarder *logging.ClusterLogForwarderSpec) sets.String {
	types := sets.NewString()
	specs := forwarder.InputMap()
	for inputName := range logging.NewRoutes(forwarder.Pipelines).ByInput {
		if logging.ReservedInputNames.Has(inputName) {
			types.Insert(inputName) // Use name as type.
		} else if spec, ok := specs[inputName]; ok {
			if spec.Application != nil {
				types.Insert(logging.InputNameApplication)
			}
			if spec.Infrastructure != nil {
				types.Insert(logging.InputNameInfrastructure)
			}
			if spec.Audit != nil {
				types.Insert(logging.InputNameAudit)
			}
		}
	}
	return types
}

func (CLO) InputsToPipelines(fwdspec *logging.ClusterLogForwarderSpec) logging.RouteMap {
	result := logging.RouteMap{}
	inputs := fwdspec.InputMap()
	for _, pipeline := range fwdspec.Pipelines {
		for _, inRef := range pipeline.InputRefs {
			if input, ok := inputs[inRef]; ok {
				// User defined input spec, unwrap.
				for t := range input.Types() {
					result.Insert(t, pipeline.Name)
				}
			} else {
				// Not a user defined type, insert direct.
				result.Insert(inRef, pipeline.Name)
			}
		}
	}
	return result
}
