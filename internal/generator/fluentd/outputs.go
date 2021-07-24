package fluentd

import (
	"fmt"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/elasticsearch"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/fluentdforward"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/kafka"
	"github.com/vimalk78/collector-conf-gen/internal/generator/helpers"
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
		var matchElement Element
		switch o.Type {
		case logging.OutputTypeElasticsearch:
			matchElement = elasticsearch.Conf(bufspec, secret, o, op)
		case logging.OutputTypeFluentdForward:
			matchElement = fluentdforward.Conf(bufspec, secret, o, op)
		case logging.OutputTypeKafka:
			matchElement = kafka.Conf(bufspec, secret, o, op)
		}
		output := FromLabel{
			Desc:    "Output to elasticsearch",
			InLabel: labelName(o.Name),
			SubElements: []Element{
				Match{
					MatchTags:    "**",
					MatchElement: matchElement,
				},
			},
		}
		need_retry := true
		if o.Type == logging.OutputTypeElasticsearch && need_retry {
			es, ok := matchElement.(elasticsearch.Elasticsearch)
			if ok {
				es.RetryTag = fmt.Sprintf("retry_%s", strings.ToLower(helpers.Replacer.Replace(o.Name)))
				output.SubElements = []Element{
					Match{
						MatchTags:    es.RetryTag,
						MatchElement: elasticsearch.RetryConf(bufspec, secret, o, op),
					},
					Match{
						MatchTags:    "**",
						MatchElement: es,
					},
				}
			}
		}
		outputs = append(outputs, output)
	}
	return outputs
}
