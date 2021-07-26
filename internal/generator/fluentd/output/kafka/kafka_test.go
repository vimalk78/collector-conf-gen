package kafka

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

var kafka_store_test = Describe("Generate fluentd config", func() {
	var f = func(clspec logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec logging.ClusterLogForwarderSpec, op Options) []Element {
		var bufspec *logging.FluentdBufferSpec = nil
		if clspec.Forwarder != nil &&
			clspec.Forwarder.Fluentd != nil &&
			clspec.Forwarder.Fluentd.Buffer != nil {
			bufspec = clspec.Forwarder.Fluentd.Buffer
		}
		return Conf(bufspec, secrets[clfspec.Outputs[0].Name], clfspec.Outputs[0], &Options{})
	}
	DescribeTable("for kafka store", TestGenerateConfWith(f),
		Entry("with username,password to single topic", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Type: logging.OutputTypeKafka,
						Name: "kafka-receiver",
						URL:  "tls://broker1-kafka.svc.messaging.cluster.local:9092/topic",
						Secret: &logging.OutputSecretSpec{
							Name: "kafka-receiver-1",
						},
						OutputTypeSpec: logging.OutputTypeSpec{
							Kafka: &logging.Kafka{
								Topic: "build_complete",
							},
						},
					},
				},
			},
			Secrets: map[string]*corev1.Secret{
				"kafka-receiver": &corev1.Secret{
					Data: map[string][]byte{
						"username": []byte("junk"),
						"password": []byte("junk"),
					},
				},
			},
			ExpectedConf: `
# Output to kafka
<label @KAFKA_RECEIVER>
  <match **>
    @type kafka2
    @id kafka_receiver
    brokers broker1-kafka.svc.messaging.cluster.local:9092
    default_topic build_complete
    use_event_time true
    sasl_plain_username "#{File.exists?('/var/run/ocp-collector/secrets/username') ? open('/var/run/ocp-collector/secrets/username','r') do |f|f.read end : ''}"
    sasl_plain_password "#{File.exists?('/var/run/ocp-collector/secrets/password') ? open('/var/run/ocp-collector/secrets/password','r') do |f|f.read end : ''}"
    sasl_over_ssl false
    <format>
      @type json
    </format>
    <buffer build_complete>
      @type file
      path '/var/lib/fluentd/kafka_receiver'
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
						Type: logging.OutputTypeKafka,
						Name: "kafka-receiver",
						URL:  "tls://broker1-kafka.svc.messaging.cluster.local:9092/topic",
						Secret: &logging.OutputSecretSpec{
							Name: "kafka-receiver-1",
						},
					},
				},
			},
			Secrets: map[string]*corev1.Secret{
				"kafka-receiver": &corev1.Secret{
					Data: map[string][]byte{
						"tls.key":       []byte("junk"),
						"tls.crt":       []byte("junk"),
						"ca-bundle.crt": []byte("junk"),
					},
				},
			},
			ExpectedConf: `
# Output to kafka
<label @KAFKA_RECEIVER>
  <match **>
    @type kafka2
    @id kafka_receiver
    brokers broker1-kafka.svc.messaging.cluster.local:9092
    default_topic topic
    use_event_time true
    client_key /var/run/ocp-collector/secrets/tls.key
    client_cert /var/run/ocp-collector/secrets/tls.crt
    ca_file /var/run/ocp-collector/secrets/ca-bundle.crt
    sasl_over_ssl false
    <format>
      @type json
    </format>
    <buffer topic>
      @type file
      path '/var/lib/fluentd/kafka_receiver'
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
						Type: logging.OutputTypeKafka,
						Name: "kafka-receiver",
						URL:  "tls://broker1-kafka.svc.messaging.cluster.local:9092/topic",
					},
				},
			},
			Secrets: security.NoSecrets,
			ExpectedConf: `
# Output to kafka
<label @KAFKA_RECEIVER>
  <match **>
    @type kafka2
    @id kafka_receiver
    brokers broker1-kafka.svc.messaging.cluster.local:9092
    default_topic topic
    use_event_time true
    <format>
      @type json
    </format>
    <buffer topic>
      @type file
      path '/var/lib/fluentd/kafka_receiver'
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
