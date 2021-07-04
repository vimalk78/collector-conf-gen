package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

var clo CLO

func MakeGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) MakeLoggingConf(spec *logging.ClusterLogForwarderSpec) []Section {
	return []Section{
		{
			g.Sources(spec),
			"Set of all input sources",
		},
		{
			g.PrometheusMetrics(),
			"Section to add measurement, and dispatch to Concat or Ingress pipelines",
		},
		{
			g.Concat(),
			`Concat pipeline 
			section`,
		},
		{
			g.Ingress(spec),
			"Ingress pipeline",
		},
		{
			g.InputsToPipeline(spec),
			"Inputs go to pipelines",
		},
		{
			g.PipelineToOutputs(spec),
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
