package vector

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

var log_sources_test = Describe("Vector Config Generation", func() {
	var f = func(clspec logging.ClusterLoggingSpec, clfspec logging.ClusterLogForwarderSpec) []Element {
		a := MakeConf()
		return MergeElements(
			a.LogSources(&clfspec, &Options{}),
		)
	}
	DescribeTable("Source(s)", TestGenerateConfWith(f),
		Entry("Only Application", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs: []string{
							logging.InputNameApplication,
						},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Logs from containers (including openshift containers)
[sources.container_logs]
  auto_partial_merge = true
  exclude_paths_glob_patterns = ["/var/log/containers/vector-*_openshift-logging_*.log", "/var/log/containers/elasticsearch-*_openshift-logging_*.log", "/var/log/containers/kibana-*_openshift-logging_*.log"]
`,
		}),
		Entry("Only Infrastructure", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs: []string{
							logging.InputNameInfrastructure,
						},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Logs from containers (including openshift containers)
[sources.container_logs]
  auto_partial_merge = true
  exclude_paths_glob_patterns = ["/var/log/containers/vector-*_openshift-logging_*.log", "/var/log/containers/elasticsearch-*_openshift-logging_*.log", "/var/log/containers/kibana-*_openshift-logging_*.log"]

[sources.journal_logs]
  type = "journald"
`,
		}),
		Entry("Only Audit", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs: []string{
							logging.InputNameAudit,
						},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Logs from host audit
[sources.host_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/audit/audit.log"]

# Logs from kubernetes audit
[sources.k8s_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/kube-apiserver/audit.log"]

# Logs from openshift audit
[sources.openshift_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/oauth-apiserver.audit.log"]
`,
		}),
		Entry("All Log Sources", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs: []string{
							logging.InputNameApplication,
							logging.InputNameInfrastructure,
							logging.InputNameAudit,
						},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Logs from containers (including openshift containers)
[sources.container_logs]
  auto_partial_merge = true
  exclude_paths_glob_patterns = ["/var/log/containers/vector-*_openshift-logging_*.log", "/var/log/containers/elasticsearch-*_openshift-logging_*.log", "/var/log/containers/kibana-*_openshift-logging_*.log"]

[sources.journal_logs]
  type = "journald"

# Logs from host audit
[sources.host_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/audit/audit.log"]

# Logs from kubernetes audit
[sources.k8s_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/kube-apiserver/audit.log"]

# Logs from openshift audit
[sources.openshift_audit_logs]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/oauth-apiserver.audit.log"]`,
		}))
})

/*
var metric_sources_test = Describe("Testing Config Generation", func() {
	var f = func(spec logging.ClusterLogForwarderSpec) []Element {
		a := MakeConf()
		return MergeElements(
			a.MetricSources(&spec, &Options{}),
		)
	}
	DescribeTable("Metric Source(s)", TestGenerateConfWith(f),
		Entry("Any Input", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Prometheus Monitoring
<source>
  @type prometheus
  bind "#{ENV['POD_IP']}"
  <ssl>
    enable true
    certificate_path "#{ENV['METRICS_CERT'] || '/etc/fluent/metrics/tls.crt'}"
    private_key_path "#{ENV['METRICS_KEY'] || '/etc/fluent/metrics/tls.key'}"
  </ssl>
</source>

<source>
  @type prometheus_monitor
  <labels>
    hostname ${hostname}
  </labels>
</source>

# excluding prometheus_tail_monitor
# since it leaks namespace/pod info
# via file paths

# This is considered experimental by the repo
<source>
  @type prometheus_output_monitor
  <labels>
    hostname ${hostname}
  </labels>
</source>`,
		}))
})
*/
