package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"k8s.io/apimachinery/pkg/util/sets"
)

var clo CLO

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

/**
 These unexported methods are copied frm CLO
**/
type CLO int

//GatherSources collects the set of unique source types and namespaces
func (CLO) GatherSources(forwarder *logging.ClusterLogForwarderSpec) sets.String {
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
