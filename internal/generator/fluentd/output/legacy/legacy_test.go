package legacy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd"
	corev1 "k8s.io/api/core/v1"
)

var legacy_test = Describe("fluentd conf generation", func() {
	Describe("generate legacy fluentdforward conf", func() {
		var f = func(clspec logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec logging.ClusterLogForwarderSpec, op Options) []Element {
			return MergeElements(
				fluentd.InputsToPipeline(&clfspec, &op),
				fluentd.PipelineToOutputs(&clfspec, &op),
				fluentd.Outputs(&clspec, secrets, &clfspec, &op))
		}
		DescribeTable("for fluentdforward store", TestGenerateConfWith(f),
			Entry("", ConfGenerateTest{
				CLFSpec: logging.ClusterLogForwarderSpec{},
				Options: Options{
					IncludeLegacyForwardConfig: "",
					IncludeLegacySyslogConfig:  "",
				},
				ExpectedConf: `
# Copying application source type to pipeline
<label @_APPLICATION>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @_LEGACY_SECUREFORWARD
    </store>
    
    <store>
      @type relabel
      @label @_LEGACY_SYSLOG
    </store>
  </match>
</label>

# Copying infrastructure source type to pipeline
<label @_INFRASTRUCTURE>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @_LEGACY_SECUREFORWARD
    </store>
    
    <store>
      @type relabel
      @label @_LEGACY_SYSLOG
    </store>
  </match>
</label>

# Copying audit source type to pipeline
<label @_AUDIT>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @_LEGACY_SECUREFORWARD
    </store>
    
    <store>
      @type relabel
      @label @_LEGACY_SYSLOG
    </store>
  </match>
</label>

<label @_LEGACY_SECUREFORWARD>  
  <match **>  
    @type copy  
    #include legacy secure-forward.conf  
    @include /etc/fluent/configs.d/secure-forward/secure-forward.conf  
  </match>  
</label>  

<label @_LEGACY_SYSLOG>
  <match **>
    @type copy
    #include legacy Syslog
    @include /etc/fluent/configs.d/syslog/syslog.conf
  </match>  
</label>
`,
			}))
	})
})

func TestFluendConfGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fluend Conf Generation")
}
