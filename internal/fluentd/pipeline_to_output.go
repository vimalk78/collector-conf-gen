package fluentd

import (
	"sort"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

func (g *Generator) PipelineToOutputs(spec *logging.ClusterLogForwarderSpec) []Element {
	pipelines := spec.Pipelines
	sort.Slice(pipelines, func(i, j int) bool {
		return pipelines[i].Name < pipelines[j].Name
	})
	for _, p := range pipelines {
		_ = p
	}
	return nil
}
