package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/elasticsearch"
	corev1 "k8s.io/api/core/v1"
)

func Outputs(clspec *logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec *logging.ClusterLogForwarderSpec, op *Options) []Element {
	outputs := []Element{}
	var bufspec *logging.FluentdBufferSpec = nil
	if clspec != nil &&
		clspec.Forwarder != nil &&
		clspec.Forwarder.Fluentd != nil &&
		clspec.Forwarder.Fluentd.Buffer != nil {
		bufspec = clspec.Forwarder.Fluentd.Buffer
	}
	for _, o := range clfspec.Outputs {
		secret := secrets[o.Name]
		switch o.Type {
		case logging.OutputTypeElasticsearch:
			outputs = MergeElements(outputs, elasticsearch.Store(bufspec, secret, o, op))
		case logging.OutputTypeFluentdForward:
		}
	}
	return outputs
}
