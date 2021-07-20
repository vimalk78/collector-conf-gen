package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/elasticsearch"
)

func Outputs(spec *logging.ClusterLogForwarderSpec, op *Options) []Element {
	outputs := []Element{}
	for _, o := range spec.Outputs {
		switch o.Type {
		case logging.OutputTypeElasticsearch:
			outputs = MergeElements(outputs, elasticsearch.Store(o, op))
		case logging.OutputTypeFluentdForward:
		}
	}
	return outputs
}
