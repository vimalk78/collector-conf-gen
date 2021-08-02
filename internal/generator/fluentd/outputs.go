package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/elasticsearch"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/fluentdforward"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/kafka"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/legacy"
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
			outputs = MergeElements(outputs, elasticsearch.Conf(bufspec, secret, o, op))
		case logging.OutputTypeFluentdForward:
			outputs = MergeElements(outputs, fluentdforward.Conf(bufspec, secret, o, op))
		case logging.OutputTypeKafka:
			outputs = MergeElements(outputs, kafka.Conf(bufspec, secret, o, op))
		}
	}
	if Clo.IncludeLegacyForwardConfig(*op) {
		outputs = append(outputs, ConfLiteral{
			TemplateName: "legacySecureForward",
			TemplateStr:  legacy.LegacySecureForwardTemplate,
		})
	}
	if Clo.IncludeLegacySyslogConfig(*op) {
		outputs = append(outputs, ConfLiteral{
			TemplateName: "legacySyslog",
			TemplateStr:  legacy.LegacySyslogForwardTemplate,
		})
	}
	return outputs
}
