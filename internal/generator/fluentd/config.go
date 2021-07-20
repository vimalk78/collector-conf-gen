package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

// keep no state in Conf
//type Conf int

//func MakeConf() Conf {
//	return Conf(0)
//}

//func (a Conf) Assemble(spec *logging.ClusterLogForwarderSpec, o *Options) []Section {
func Conf(spec *logging.ClusterLogForwarderSpec, o *Options) []Section {
	return []Section{
		{
			Sources(spec, o),
			"Set of all input sources",
		},
		{
			PrometheusMetrics(spec, o),
			"Section to add measurement, and dispatch to Concat or Ingress pipelines",
		},
		{
			Concat(spec, o),
			`Concat pipeline 
			section`,
		},
		{
			Ingress(spec, o),
			"Ingress pipeline",
		},
		// input ends
		// give a hook here
		{
			InputsToPipeline(spec, o),
			"Inputs go to pipelines",
		},
		{
			PipelineToOutputs(spec, o),
			"Pipeline to Outputs",
		},
		// output begins here
		// give a hook here
		{
			Outputs(spec, o),
			"Outputs",
		},
	}
}
