package assembler

import (
	"fmt"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

func (a Assembler) MetricSources(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return []Element{
		PrometheusMonitor{
			Desc:         "Prometheus Monitoring",
			TemplateName: "PrometheusMonitor",
			TemplateStr:  PrometheusMonitorTemplate,
		},
	}
}

func (a Assembler) Sources(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return MergeElements(
		a.MetricSources(spec, o),
		a.LogSources(spec, o),
	)
}

//TODO: handle the following options here
// - includeLegacyForwardConfig
// - includeLegacySyslogConfig
func (a Assembler) LogSources(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	var el []Element = make([]Element, 0)
	types := clo.GatherSources(spec)
	if types.Has(logging.InputNameApplication) || types.Has(logging.InputNameInfrastructure) {
		el = append(el,
			ContainerLogs{
				Desc:         "Logs from containers (including openshift containers)",
				Paths:        ContainerLogPaths(),
				ExcludePaths: ExcludeContainerPaths(),
				PosFile:      "/var/log/es-containers.log.pos",
				OutLabel:     "MEASURE",
			})
	}
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el,
			JournalLog{
				Desc:         "Logs from linux journal",
				OutLabel:     "MEASURE",
				TemplateName: "inputSourceJournalTemplate",
				TemplateStr:  JournalLogTemplate,
			})
	}
	if types.Has(logging.InputNameAudit) {
		el = append(el,
			HostAuditLog{
				Desc:         "Logs from host audit",
				OutLabel:     "MEASURE",
				TemplateName: "inputSourceHostAuditTemplate",
				TemplateStr:  HostAuditLogTemplate,
			},
			K8sAuditLog{
				Desc:         "Logs from kubernetes audit",
				OutLabel:     "MEASURE",
				TemplateName: "inputSourceK8sAuditTemplate",
				TemplateStr:  K8sAuditLogTemplate,
			},
			OpenshiftAuditLog{
				Desc:         "Logs from openshift audit",
				OutLabel:     "MEASURE",
				TemplateName: "inputSourceOpenShiftAuditTemplate",
				TemplateStr:  OpenshiftAuditLogTemplate,
			})
	}
	return el
}

func ContainerLogPaths() string {
	return fmt.Sprintf("%q", "/var/log/containers/*.log")
}

func ExcludeContainerPaths() string {
	return fmt.Sprintf("[%s]", strings.Join(
		[]string{
			fmt.Sprintf("%q", fmt.Sprintf(CollectorLogsPath(), FluentdCollectorPodNamePrefix())),
			//fmt.Sprintf("%q", fmt.Sprintf(CollectorLogsPath(), FluentBitCollectorPodNamePrefix())),
			fmt.Sprintf("%q", fmt.Sprintf(LogStoreLogsPath(), ESLogStorePodNamePrefix())),
			//fmt.Sprintf("%q", fmt.Sprintf(LogStoreLogsPath(), LokiLogStorePodNamePrefix())),
			fmt.Sprintf("%q", VisualizationLogsPath()),
		},
		",",
	))
}

func LoggingNamespace() string {
	return "openshift-logging"
}

func FluentdCollectorPodNamePrefix() string {
	return "fluentd"
}

func FluentBitCollectorPodNamePrefix() string {
	return "fluent-bit"
}

func VectorCollectorPodNamePrefix() string {
	return "vector"
}

func CollectorLogsPath() string {
	return fmt.Sprintf("/var/log/containers/%%s-*_%s_*.log", LoggingNamespace())
}

func ESLogStorePodNamePrefix() string {
	return "elasticsearch"
}

func LokiLogStorePodNamePrefix() string {
	return "loki"
}
func LogStoreLogsPath() string {
	return fmt.Sprintf("/var/log/containers/%%s-*_%s_*.log", LoggingNamespace())
}

func VisualizationPodNamePrefix() string {
	return "kibana"
}

func VisualizationLogsPath() string {
	return fmt.Sprintf("/var/log/containers/%s-*_%s_*.log", VisualizationPodNamePrefix(), LoggingNamespace())
}
