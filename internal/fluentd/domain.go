package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

var Logging []Section = []Section{
	{
		Sources,
		"Set of all input sources",
	},
	{
		Metrics,
		"Section to add measurement, and dispatch to Concat or Ingress pipelines",
	},
	{
		Concat,
		"Concat pipeline",
	},
	{
		Ingress,
		"Ingress pipeline",
	},
}

var Sources []Element = []Element{
	PrometheusMonitor{
		Desc: "Prometheus Monitoring",
	},
	ContainerLogs{
		OutLabel: "MEASURE",
		Desc:     "Logs from containers (including openshift containers)",
	},
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
	},
	JournalLog{
		OutLabel: "MEASURE",
		Desc:     "Logs from linux journal",
	},
}

var Metrics []Element = []Element{
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

var Concat []Element = []Element{
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
			Stdout{
				Desc: "Print to standard output",
			},
		},
	},
}

var Ingress []Element = []Element{
	Pipeline{
		InLabel: "INGRESS",
		Desc:    "Concat log lines of icontainer logs",
		SubElements: []Element{
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
				TemplateName: "retagJournal",
				TemplateStr:  RetagJournalLogs,
				OutLabel:     "INGRESS",
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
	},
}

//gatherSources collects the set of unique source types and namespaces
func gatherSources(forwarder *logging.ClusterLogForwarderSpec) sets.String {
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

func SelectLogType(spec *logging.ClusterLogForwarderSpec) Section {
	var el []Element = make([]Element, 0)
	types := gatherSources(spec)
	if types.Has(logging.InputNameApplication) {
		el = append(el, Relabel{
			Desc:     "Dont discard Application logs",
			Pattern:  "kubernetes.**",
			OutLabel: "_APPLICATION",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Application logs",
			Pattern:      "kubernetes.**",
			TemplateStr:  DiscardMatched,
			TemplateName: "discardMatched",
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
			TemplateStr:  DiscardMatched,
			TemplateName: "discardMatched",
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
			TemplateStr:  DiscardMatched,
			TemplateName: "discardMatched",
		})
	}
	return Section{
		Elements: el,
		Comment:  "Generated section to send log types to its pipelines",
	}
}
