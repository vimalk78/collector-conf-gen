package elasticsearch

import (
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
)

type Elasticsearch struct {
	StoreID        string
	Host           string
	Port           string
	SecurityConfig []Element
	BufferConfig   []Element
}

func (e Elasticsearch) Name() string {
	return "elasticsearchTemplate"
}

func (e Elasticsearch) Template() string {
	return `{{define "` + e.Name() + `" -}}
<store>
  @type elasticsearch
  @id {{.StoreID}}
  host {{.Host}}
  port {{.Port}}
{{compose .SecurityConfig | indent 2}}
  verify_es_version_at_startup false
  target_index_key viaq_index_name
  id_key viaq_msg_id
  remove_keys viaq_index_name
  type_name _doc
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
{{compose .BufferConfig | indent 2}}
</store>
{{- end}}
`
}

func (e Elasticsearch) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(e.Template()))
}

func (e Elasticsearch) Data() interface{} {
	return e
}

func Store(o logging.OutputSpec, op *Options) []Element {
	return []Element{
		Elasticsearch{
			BufferConfig: output.Buffer([]string{}, nil, &o),
		},
	}
}
