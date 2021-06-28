package fluentd

import logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"

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
				g.SourceToInput(spec)),
		},
	}
}
