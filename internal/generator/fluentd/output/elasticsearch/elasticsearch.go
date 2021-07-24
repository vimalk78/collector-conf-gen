package elasticsearch

import (
	"fmt"
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	"github.com/vimalk78/collector-conf-gen/internal/generator/helpers"
	"github.com/vimalk78/collector-conf-gen/internal/generator/url"
	urlhelper "github.com/vimalk78/collector-conf-gen/internal/generator/url"
)

const (
	defaultElasticsearchPort = "9200"
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
# {{.Desc}}
@type elasticsearch
@id {{.StoreID}}
host {{.Host}}
port {{.Port}}
{{compose .SecurityConfig}}
verify_es_version_at_startup false
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

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Elasticsearch {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	port := u.Port()
	if port == "" {
		port = defaultElasticsearchPort
	}
	return Elasticsearch{
		Desc:           "Elasticsearch store",
		StoreID:        strings.ToLower(fmt.Sprintf("%v", helpers.Replacer.Replace(o.Name))),
		Host:           u.Hostname(),
		Port:           port,
		SecurityConfig: SecurityConfig(o, secret),
		BufferConfig:   output.Buffer(output.NOKEYS, bufspec, &o),
	}
}

func RetryConf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Elasticsearch {
	es := Conf(bufspec, secret, o, op)
	es.StoreID = strings.ToLower(fmt.Sprintf("retry_%v", helpers.Replacer.Replace(o.Name)))
	return es
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	tls := TLS(url.IsTLSScheme(u.Scheme) || secret != nil)
	conf := []Element{
		tls,
	}
	if security.HasUsernamePassword(secret) {
		up := UserNamePass{
			// TODO: use constants.ClientUsername
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
	return conf
}
