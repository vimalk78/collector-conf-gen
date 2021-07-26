package elasticsearch

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	corev1 "k8s.io/api/core/v1"
)

var elasticsearch_store_test = Describe("Generate fluentd config", func() {
	var f = func(clspec logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec logging.ClusterLogForwarderSpec, op Options) []Element {
		var bufspec *logging.FluentdBufferSpec = nil
		if clspec.Forwarder != nil &&
			clspec.Forwarder.Fluentd != nil &&
			clspec.Forwarder.Fluentd.Buffer != nil {
			bufspec = clspec.Forwarder.Fluentd.Buffer
		}
		return Conf(bufspec, secrets[clfspec.Outputs[0].Name], clfspec.Outputs[0], &Options{})
	}
	DescribeTable("for Elasticsearch store", TestGenerateConfWith(f),
		Entry("with username,password", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Type: logging.OutputTypeElasticsearch,
						Name: "es-1",
						URL:  "https://es.svc.infra.cluster:9999",
						Secret: &logging.OutputSecretSpec{
							Name: "es-1",
						},
					},
				},
			},
			Secrets: map[string]*corev1.Secret{
				"es-1": &corev1.Secret{
					Data: map[string][]byte{
						"username": []byte("junk"),
						"password": []byte("junk"),
					},
				},
			},
			ExpectedConf: `
# Output to elasticsearch
<label @ES_1>
  #remove structured field if present
  <filter **>
    @type record_modifier
    remove_keys structured
  </filter>
  
  #flatten labels to prevent field explosion in ES
  <filter **>
    @type record_modifier
    <record>
      kubernetes ${!record['kubernetes'].nil? ? record['kubernetes'].merge({"flat_labels": (record['kubernetes']['labels']||{}).map{|k,v| "#{k}=#{v}"}}) : {} }
    </record>
    remove_keys $.kubernetes.labels
  </filter>
  
  <match retry_es_1>
    # Elasticsearch store
    @type elasticsearch
    @id retry_es_1
    host es.svc.infra.cluster
    port 9999
    scheme https
    ssl_version TLSv1_2
    username "#{File.exists?('/var/run/ocp-collector/secrets/username') ? open('/var/run/ocp-collector/secrets/username','r') do |f|f.read end : ''}"
    password "#{File.exists?('/var/run/ocp-collector/secrets/password') ? open('/var/run/ocp-collector/secrets/password','r') do |f|f.read end : ''}"
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
  
  <match **>
    # Elasticsearch store
    @type elasticsearch
    @id es_1
    host es.svc.infra.cluster
    port 9999
    scheme https
    ssl_version TLSv1_2
    username "#{File.exists?('/var/run/ocp-collector/secrets/username') ? open('/var/run/ocp-collector/secrets/username','r') do |f|f.read end : ''}"
    password "#{File.exists?('/var/run/ocp-collector/secrets/password') ? open('/var/run/ocp-collector/secrets/password','r') do |f|f.read end : ''}"
    verify_es_version_at_startup false
    target_index_key viaq_index_name
    id_key viaq_msg_id
    remove_keys viaq_index_name
    type_name _doc
    retry_tag retry_es_1
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
</label>
`,
		}),
		Entry("with tls key,cert,ca-bundle", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Type: logging.OutputTypeElasticsearch,
						Name: "es-1",
						URL:  "https://es.svc.infra.cluster:9999",
						Secret: &logging.OutputSecretSpec{
							Name: "es-1",
						},
					},
				},
			},
			Secrets: map[string]*corev1.Secret{
				"es-1": &corev1.Secret{
					Data: map[string][]byte{
						"tls.key":       []byte("junk"),
						"tls.crt":       []byte("junk"),
						"ca-bundle.crt": []byte("junk"),
					},
				},
			},
			ExpectedConf: `
# Output to elasticsearch
<label @ES_1>
  #remove structured field if present
  <filter **>
    @type record_modifier
    remove_keys structured
  </filter>
  
  #flatten labels to prevent field explosion in ES
  <filter **>
    @type record_modifier
    <record>
      kubernetes ${!record['kubernetes'].nil? ? record['kubernetes'].merge({"flat_labels": (record['kubernetes']['labels']||{}).map{|k,v| "#{k}=#{v}"}}) : {} }
    </record>
    remove_keys $.kubernetes.labels
  </filter>
  
  <match retry_es_1>
    # Elasticsearch store
    @type elasticsearch
    @id retry_es_1
    host es.svc.infra.cluster
    port 9999
    scheme https
    ssl_version TLSv1_2
    client_key /var/run/ocp-collector/secrets/tls.key
    client_cert /var/run/ocp-collector/secrets/tls.crt
    ca_file /var/run/ocp-collector/secrets/ca-bundle.crt
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
  
  <match **>
    # Elasticsearch store
    @type elasticsearch
    @id es_1
    host es.svc.infra.cluster
    port 9999
    scheme https
    ssl_version TLSv1_2
    client_key /var/run/ocp-collector/secrets/tls.key
    client_cert /var/run/ocp-collector/secrets/tls.crt
    ca_file /var/run/ocp-collector/secrets/ca-bundle.crt
    verify_es_version_at_startup false
    target_index_key viaq_index_name
    id_key viaq_msg_id
    remove_keys viaq_index_name
    type_name _doc
    retry_tag retry_es_1
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
</label>
`,
		}),
		Entry("without security", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Type: logging.OutputTypeElasticsearch,
						Name: "es-1",
						URL:  "http://es.svc.infra.cluster:9999",
						Secret: &logging.OutputSecretSpec{
							Name: "es-1",
						},
					},
				},
			},
			Secrets: security.NoSecrets,
			ExpectedConf: `
# Output to elasticsearch
<label @ES_1>
  #remove structured field if present
  <filter **>
    @type record_modifier
    remove_keys structured
  </filter>
  
  #flatten labels to prevent field explosion in ES
  <filter **>
    @type record_modifier
    <record>
      kubernetes ${!record['kubernetes'].nil? ? record['kubernetes'].merge({"flat_labels": (record['kubernetes']['labels']||{}).map{|k,v| "#{k}=#{v}"}}) : {} }
    </record>
    remove_keys $.kubernetes.labels
  </filter>
  
  <match retry_es_1>
    # Elasticsearch store
    @type elasticsearch
    @id retry_es_1
    host es.svc.infra.cluster
    port 9999
    scheme http
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
  
  <match **>
    # Elasticsearch store
    @type elasticsearch
    @id es_1
    host es.svc.infra.cluster
    port 9999
    scheme http
    verify_es_version_at_startup false
    target_index_key viaq_index_name
    id_key viaq_msg_id
    remove_keys viaq_index_name
    type_name _doc
    retry_tag retry_es_1
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
    <buffer>
      @type file
      path '/var/lib/fluentd/es_1'
      flush_mode interval
      flush_interval 1s
      flush_thread_count 2
      flush_at_shutdown true
      retry_type exponential_backoff
      retry_wait 1s
      retry_max_interval 60s
      retry_timeout 60m
      queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
      total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] || '8589934592'}"
      chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '1m'}"
      overflow_action block
    </buffer>
  </match>
</label>
`,
		}),
	)
})

func TestFluendConfGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fluend Conf Generation")
}
