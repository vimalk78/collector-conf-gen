package kafka

import (
	"net/url"
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	"github.com/vimalk78/collector-conf-gen/internal/generator/helpers"
	urlhelper "github.com/vimalk78/collector-conf-gen/internal/generator/url"
)

const (
	defaultKafkaTopic = "topic"
)

type Kafka struct {
	Desc           string
	StoreID        string
	Brokers        string
	Topics         string
	SecurityConfig []Element
	BufferConfig   []Element
}

func (k Kafka) Name() string {
	return "kafkaTemplate"
}

func (k Kafka) Template() string {
	return `{{define "` + k.Name() + `" -}}
@type kafka2
@id {{.StoreID}}
brokers {{.Brokers}}
default_topic {{.Topics}}
use_event_time true
{{- with $x := compose .SecurityConfig }}
{{$x}}
{{- end}}
<format>
  @type json
</format>
{{compose .BufferConfig}}
{{- end}}
`
}

func (k Kafka) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(k.Template()))
}

func (k Kafka) Data() interface{} {
	return k
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Element {
	topics := Topics(o)
	return Kafka{
		Desc:           "Kafka store",
		StoreID:        strings.ToLower(helpers.Replacer.Replace(o.Name)),
		Topics:         topics,
		Brokers:        Brokers(o),
		SecurityConfig: SecurityConfig(o, secret),
		BufferConfig:   output.Buffer([]string{topics}, bufspec, &o),
	}
}

//Brokers returns the list of broker endpoints of a kafka cluster.
//The list represents only the initial set used by the collector's kafka client for the
//first connention only. The collector's kafka client fetches constantly an updated list
//from kafka. These updates are not reconciled back to the collector configuration.
//The list of brokers are populated from the Kafka OutputSpec `Brokers` field, a list of
//valid URLs. If none provided the target URL from the OutputSpec is used as fallback.
//Finally, if neither approach works the current collector process will be terminated.
func Brokers(o logging.OutputSpec) string {
	parseBroker := func(b string) string {
		url, _ := url.Parse(b)
		return url.Host
	}

	if o.Kafka != nil {
		if o.Kafka.Brokers != nil {
			brokers := []string{}
			for _, broker := range o.Kafka.Brokers {
				b := parseBroker(broker)
				if b != "" {
					brokers = append(brokers, b)
				}
			}

			if len(brokers) > 0 {
				return strings.Join(brokers, ",")
			}
		}
	}

	// Fallback to parse a single broker from target's URL
	fallback := parseBroker(o.URL)
	if fallback == "" {
		panic("Failed to parse any Kafka broker from output spec")
	}

	return fallback
}

//Topic returns the name of an existing kafka topic.
//The kafka topic is either extracted from the kafka OutputSpec `Topic` field in a multiple broker
//setup or as a fallback from the OutputSpec URL if provided as a host path. Defaults to `topic`.
func Topics(o logging.OutputSpec) string {
	if o.Kafka != nil && o.Kafka.Topic != "" {
		return o.Kafka.Topic
	}

	url, _ := urlhelper.Parse(o.URL)
	topic := strings.TrimLeft(url.Path, "/")
	if topic != "" {
		return topic
	}

	// Fallback to default topic
	return defaultKafkaTopic
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	conf := []Element{}
	if secret != nil {
		if security.HasUsernamePassword(secret) {
			up := UserNamePass{
				UsernamePath: security.SecretPath(secret, "username"),
				PasswordPath: security.SecretPath(secret, "password"),
			}
			conf = append(conf, up)
		}
		if security.HasTLSKeyAndCrt(secret) {
			kc := TLSKeyCert{
				// TODO: use constants.ClientCertKey
				KeyPath:  security.SecretPath(secret, "tls.key"),
				CertPath: security.SecretPath(secret, "tls.crt"),
			}
			conf = append(conf, kc)
		}
		if security.HasCABundle(secret) {
			ca := CAFile{
				// TODO: use constants.TrustedCABundleKey
				CAFilePath: security.SecretPath(secret, "ca-bundle.crt"),
			}
			conf = append(conf, ca)
		}
		if _, ok := secret.Data["sasl_over_ssl"]; ok {
			s := SaslOverSSL(true)
			conf = append(conf, s)
		} else {
			s := SaslOverSSL(false)
			conf = append(conf, s)
		}
	}
	return conf
}
