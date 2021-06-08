package main

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/fluentd"
	. "github.com/vimalk78/collector-conf-gen/internal/fluentd"
)

func test() {
	e := []Element{
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
			Desc:     "Logs from kubernetes audit ",
		},
		OpenshiftAuditLog{
			OutLabel: "MEASURE",
			Desc:     "Logs from openshift audit ",
		},
		JournalLog{
			OutLabel: "MEASURE",
			Desc:     "Logs from linux journal",
		},
		Pipeline{
			InLabel: "MEASURE",
			Desc:    "Handle Measure",
			SubElements: []Element{
				ConfLiteral{
					TemplateName: "HandleMeasure",
					Desc:         "xx",
					TemplateStr:  HandleMeasure,
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
		Pipeline{
			InLabel: "CONCAT",
			Desc:    "Concat log lines of icontainer logs",
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
			},
		},
	}
	conf, err := fluentd.GenerateConf(e...)
	if err != nil {
		fmt.Printf("error occured %v\n", err)
	} else {
		fmt.Printf("%s\n", conf)
	}
}

func main() {
	test()
}

func useCLF() {
	spec := logging.ClusterLogForwarderSpec{
		Inputs: []logging.InputSpec{},
	}
	fmt.Printf("spec: %#v\n", spec)
}
