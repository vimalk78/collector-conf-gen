package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	corev1 "k8s.io/api/core/v1"
)

var buffer_test = Describe("Generate fluentd conf", func() {
	var f = func(clspec logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec logging.ClusterLogForwarderSpec, op Options) []Element {
		es := make([][]Element, len(clfspec.Outputs))
		for i := range clfspec.Outputs {
			es[i] = output.Buffer([]string{"time", "tag"}, clspec.Forwarder.Fluentd.Buffer, &clfspec.Outputs[i])
		}
		return MergeElements(es...)
	}
	DescribeTable("Buffers", TestGenerateConfWith(f),
		Entry("With tuning parameters", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Name: "es-1",
					},
				},
			},
			CLSpec: logging.ClusterLoggingSpec{
				Forwarder: &logging.ForwarderSpec{
					Fluentd: &logging.FluentdForwarderSpec{
						Buffer: &logging.FluentdBufferSpec{
							ChunkLimitSize:   "8m",
							TotalLimitSize:   "800000000",
							OverflowAction:   "throw_exception",
							FlushThreadCount: 128,
							FlushMode:        "immediate",
							FlushInterval:    "25s",
							RetryWait:        "20s",
							RetryType:        "periodic",
							RetryMaxInterval: "300s",
							RetryTimeout:     "60h",
						},
					},
				},
			},
			ExpectedConf: `
<buffer time,tag>
  @type file
  path '/var/lib/fluentd/es_1'
  flush_mode immediate
  flush_interval 25s
  flush_thread_count 128
  flush_at_shutdown true
  retry_type periodic
  retry_wait 20s
  retry_max_interval 300s
  retry_timeout 60h
  queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32'}"
  total_limit_size 800000000
  chunk_limit_size 8m
  overflow_action throw_exception
</buffer>`,
		}))
})

var retry_buffer_test = Describe("", func() {
	var f = func(clspec logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec logging.ClusterLogForwarderSpec, op Options) []Element {
		es := make([][]Element, len(clfspec.Outputs))
		for i := range clfspec.Outputs {
			es[i] = output.RetryBuffer([]string{}, clspec.Forwarder.Fluentd.Buffer, &clfspec.Outputs[i])
		}
		return MergeElements(es...)
	}
	DescribeTable("Buffers", TestGenerateConfWith(f),
		Entry("With no tuning parameters", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "defaultoutput",
					},
				},
				Outputs: []logging.OutputSpec{
					{
						Name: "es-2",
					},
				},
			},
			CLSpec: logging.ClusterLoggingSpec{
				Forwarder: &logging.ForwarderSpec{
					Fluentd: &logging.FluentdForwarderSpec{
						Buffer: nil,
					},
				},
			},
			ExpectedConf: `
<buffer>
  @type file
  path '/var/lib/fluentd/retry_es_2'
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
</buffer>`,
		}))
})
