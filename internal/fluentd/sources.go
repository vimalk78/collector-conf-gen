package fluentd

import logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"

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
