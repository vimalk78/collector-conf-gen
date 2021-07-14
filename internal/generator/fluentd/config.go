package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

// keep no state in Conf
type Conf int

func MakeConf() Conf {
	return Conf(0)
}

type Options map[string]string

func (a Conf) Assemble(spec *logging.ClusterLogForwarderSpec) []Section {
	return a.AssembleConfWithOptions(spec, &Options{})
}

func (a Conf) AssembleConfWithOptions(spec *logging.ClusterLogForwarderSpec, o *Options) []Section {
	return []Section{
		{
			a.Sources(spec, o),
			"Set of all input sources",
		},
		{
			a.PrometheusMetrics(spec, o),
			"Section to add measurement, and dispatch to Concat or Ingress pipelines",
		},
		{
			a.Concat(spec, o),
			`Concat pipeline 
			section`,
		},
		{
			a.Ingress(spec, o),
			"Ingress pipeline",
		},
		{
			a.InputsToPipeline(spec, o),
			"Inputs go to pipelines",
		},
		{
			a.PipelineToOutputs(spec, o),
			"Pipeline to Outputs",
		},
	}
}
