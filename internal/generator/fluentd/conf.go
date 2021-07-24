package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	corev1 "k8s.io/api/core/v1"
)

var Header = `
## CLO GENERATED CONFIGURATION ###
# This file is a copy of the fluentd configuration entrypoint
# which should normally be supplied in a configmap.

<system>
  log_level "#{ENV['LOG_LEVEL'] || 'warn'}"
</system>

`

func Conf(clspec *logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec *logging.ClusterLogForwarderSpec, op *Options) []Section {
	return []Section{
		{
			Sources(clfspec, op),
			"Set of all input sources",
		},
		{
			PrometheusMetrics(clfspec, op),
			"Section to add measurement, and dispatch to Concat or Ingress pipelines",
		},
		{
			Concat(clfspec, op),
			`Concat pipeline 
			section`,
		},
		{
			Ingress(clfspec, op),
			"Ingress pipeline",
		},
		// input ends
		// give a hook here
		{
			InputsToPipeline(clfspec, op),
			"Inputs go to pipelines",
		},
		{
			PipelineToOutputs(clfspec, op),
			"Pipeline to Outputs",
		},
		// output begins here
		// give a hook here
		{
			Outputs(clspec, secrets, clfspec, op),
			"Outputs",
		},
	}
}
