package elasticsearch

import (
	"fmt"
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	"github.com/vimalk78/collector-conf-gen/internal/generator/url"
	urlhelper "github.com/vimalk78/collector-conf-gen/internal/generator/url"
)

const (
	defaultElasticsearchPort = "9200"
	KeyStructured            = "structured"
)

type Elasticsearch struct {
	Desc           string
	StoreID        string
	Host           string
	Port           string
	RetryTag       string
	SecurityConfig []Element
	BufferConfig   []Element
}

func (e Elasticsearch) Name() string {
	return "elasticsearchTemplate"
}

func (e Elasticsearch) Template() string {
	return `{{define "` + e.Name() + `" -}}
{{if .Desc -}}
# {{.Desc}}
{{ end -}}
@type elasticsearch
@id {{.StoreID}}
host {{.Host}}
port {{.Port}}
verify_es_version_at_startup false
{{compose .SecurityConfig}}
target_index_key viaq_index_name
id_key viaq_msg_id
remove_keys viaq_index_name
type_name _doc
{{- if .RetryTag}}
retry_tag {{.RetryTag}}
{{- end}}
http_backend typhoeus
write_operation create
reload_connections 'true'
# https://github.com/uken/fluent-plugin-elasticsearch#reload-after
reload_after '200'
# https://github.com/uken/fluent-plugin-elasticsearch#sniffer-class-name
sniffer_class_name 'Fluent::Plugin::ElasticsearchSimpleSniffer'
reload_on_failure false
# 2 ^ 31
request_timeout 2147483648
{{compose .BufferConfig}}
{{- end}}
`
}

func (e Elasticsearch) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(e.Template()))
}

func (e Elasticsearch) Data() interface{} {
	return e
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	return []Element{
		FromLabel{
			InLabel: helpers.LabelName(o.Name),
			SubElements: MergeElements(
				ChangeESIndex(bufspec, secret, o, op),
				FlattenLabels(bufspec, secret, o, op),
				OutputConf(bufspec, secret, o, op),
			),
		},
	}
}

func OutputConf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	es := ESOutput(bufspec, secret, o, op)
	need_retry := true
	if need_retry {
		es.RetryTag = helpers.StoreID(o.Name, true)
		return []Element{
			Match{
				MatchTags:    es.RetryTag,
				MatchElement: RetryESOutput(bufspec, secret, o, op),
			},
			Match{
				MatchTags:    "**",
				MatchElement: es,
			},
		}
	} else {
		return []Element{
			Match{
				MatchTags:    "**",
				MatchElement: es,
			},
		}
	}
}

func ESOutput(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Elasticsearch {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	port := u.Port()
	if port == "" {
		port = defaultElasticsearchPort
	}
	storeID := helpers.StoreID(o.Name, false)
	return Elasticsearch{
		StoreID:        storeID,
		Host:           u.Hostname(),
		Port:           port,
		SecurityConfig: SecurityConfig(o, secret),
		BufferConfig:   output.Buffer(output.NOKEYS, bufspec, storeID, &o),
	}
}

func RetryESOutput(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Elasticsearch {
	es := ESOutput(bufspec, secret, o, op)
	es.StoreID = helpers.StoreID(o.Name, true)
	es.BufferConfig = output.Buffer(output.NOKEYS, bufspec, es.StoreID, &o)
	return es
}

func ChangeESIndex(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	if o.Elasticsearch != nil && (o.Elasticsearch.StructuredTypeKey != "" || o.Elasticsearch.StructuredTypeName != "") {
		return []Element{
			Filter{
				MatchTags: "**",
				Element: RecordModifier{
					Record: map[RecordKey]RubyExpression{
						"typeFromKey":           RubyExpression(fmt.Sprintf("${record.dig(%s)}", generateRubyDigArgs(o.Elasticsearch.StructuredTypeKey))),
						"hasStructuredTypeName": RubyExpression(o.Elasticsearch.StructuredTypeName),
						"viaq_index_name":       "",
					},
					RemoveKeys: []string{"typeFromKey", "hasStructuredTypeName", "viaq_index_name"},
				},
			},
		}
	} else {
		return []Element{
			Filter{
				Desc:      "remove structured field if present",
				MatchTags: "**",
				Element: RecordModifier{
					RemoveKeys: []string{KeyStructured},
				},
			},
		}
	}
}

func FlattenLabels(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	return []Element{
		Filter{
			Desc:      "flatten labels to prevent field explosion in ES",
			MatchTags: "**",
			Element: RecordModifier{
				Record: map[RecordKey]RubyExpression{
					"kubernetes": `${!record['kubernetes'].nil? ? record['kubernetes'].merge({"flat_labels": (record['kubernetes']['labels']||{}).map{|k,v| "#{k}=#{v}"}}) : {} }`,
				},
				RemoveKeys: []string{"$.kubernetes.labels"},
			},
		},
	}
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	tls := TLS(url.IsTLSScheme(u.Scheme) || secret != nil)
	conf := []Element{
		tls,
	}
	if o.Secret == nil {
		return conf
	}
	if security.HasUsernamePassword(secret) {
		up := UserNamePass{
			// TODO: use constants.ClientUsername
			UsernamePath: security.SecretPath(o.Secret.Name, "username"),
			PasswordPath: security.SecretPath(o.Secret.Name, "password"),
		}
		conf = append(conf, up)
	}
	if security.HasTLSKeyAndCrt(secret) {
		kc := TLSKeyCert{
			// TODO: use constants.ClientCertKey
			KeyPath:  security.SecretPath(o.Secret.Name, "tls.key"),
			CertPath: security.SecretPath(o.Secret.Name, "tls.crt"),
		}
		conf = append(conf, kc)
	}
	if security.HasCABundle(secret) {
		ca := CAFile{
			// TODO: use constants.TrustedCABundleKey
			CAFilePath: security.SecretPath(o.Secret.Name, "ca-bundle.crt"),
		}
		conf = append(conf, ca)
	}
	return conf
}

func generateRubyDigArgs(path string) string {
	var args []string
	for _, s := range strings.Split(path, ".") {
		args = append(args, fmt.Sprintf("%q", s))
	}
	return strings.Join(args, ",")
}
