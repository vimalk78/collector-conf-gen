package vector

import (
	"fmt"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/vector/source"
)

func (a Conf) Sources(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return MergeElements(
		a.LogSources(spec, o),
	)
}

func (a Conf) LogSources(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	var el []Element = make([]Element, 0)
	types := Clo.GatherSources(spec, o)
	if types.Has(logging.InputNameApplication) || types.Has(logging.InputNameInfrastructure) {
		el = append(el,
			source.KubernetesLogs{
				ComponentID:  "container_logs",
				Desc:         "Logs from containers (including openshift containers)",
				ExcludePaths: ExcludeContainerPaths(),
			})
	}
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el,
			source.JournalLog{
				ComponentID:  "journal_logs",
				Desc:         "Logs from linux journal",
				TemplateName: "inputSourceJournalTemplate",
				TemplateStr:  source.JournalLogTemplate,
			})
	}
	if types.Has(logging.InputNameAudit) {
		el = append(el,
			source.HostAuditLog{
				ComponentID:  "host_audit_logs",
				Desc:         "Logs from host audit",
				TemplateName: "inputSourceHostAuditTemplate",
				TemplateStr:  source.HostAuditLogTemplate,
			},
			source.K8sAuditLog{
				ComponentID:  "k8s_audit_logs",
				Desc:         "Logs from kubernetes audit",
				TemplateName: "inputSourceK8sAuditTemplate",
				TemplateStr:  source.K8sAuditLogTemplate,
			},
			source.OpenshiftAuditLog{
				ComponentID:  "openshift_audit_logs",
				Desc:         "Logs from openshift audit",
				TemplateName: "inputSourceOpenShiftAuditTemplate",
				TemplateStr:  source.OpenshiftAuditLogTemplate,
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
			fmt.Sprintf("%q", fmt.Sprintf(CollectorLogsPath(), VectorCollectorPodNamePrefix())),
			fmt.Sprintf("%q", fmt.Sprintf(LogStoreLogsPath(), ESLogStorePodNamePrefix())),
			fmt.Sprintf("%q", VisualizationLogsPath()),
		},
		", ",
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
