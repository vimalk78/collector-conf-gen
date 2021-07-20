package fluentd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

var pipeline_to_outputs = Describe("Testing Config Generation", func() {
	var f = func(clspec logging.ClusterLoggingSpec, clfspec logging.ClusterLogForwarderSpec) []Element {
		//a := MakeConf()
		return MergeElements(
			PipelineToOutputs(&clfspec, &Options{}),
		)
	}
	DescribeTable("Pipelines(s) to Output(s)", TestGenerateConfWith(f),
		Entry("Application to single output", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "defaultoutput",
					},
				},
			},
			ExpectedConf: `
# Copying pipeline defaultoutput to outputs
<label @DEFAULTOUTPUT>
  <match **>
    @type relabel
    @label @DEFAULT
  </match>
</label>`,
		}),
		Entry("Application to multiple outputs", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault, "es-app-out"},
						Name:       "app-to-es",
					},
					{
						InputRefs:  []string{logging.InputNameAudit},
						OutputRefs: []string{logging.OutputNameDefault, "es-audit-out"},
						Name:       "audit-to-es",
					},
				},
			},
			ExpectedConf: `
# Copying pipeline app-to-es to outputs
<label @APP_TO_ES>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @DEFAULT
    </store>
    <store>
      @type relabel
      @label @ES_APP_OUT
    </store>
  </match>
</label>
# Copying pipeline audit-to-es to outputs
<label @AUDIT_TO_ES>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @DEFAULT
    </store>
    <store>
      @type relabel
      @label @ES_AUDIT_OUT
    </store>
  </match>
</label>`,
		}),
		Entry("Application to default output with Labels", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault, "es-app-out"},
						Name:       "app-to-es",
						Labels: map[string]string{
							"a": "b",
							"c": "d",
						},
					},
				},
			},
			ExpectedConf: `
# Copying pipeline app-to-es to outputs
<label @APP_TO_ES>
  # Add User Defined labels to the output record
  <filter **>
    @type record_transformer
    <record>
      openshift { "labels": {"a":"b","c":"d"} }
    </record>
  </filter>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @DEFAULT
    </store>
    <store>
      @type relabel
      @label @ES_APP_OUT
    </store>
  </match>
</label>`,
		}),
		Entry("Application to default output with Json Parsing", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault, "es-app-out"},
						Name:       "app-to-es",
						Parse:      "json",
					},
				},
			},
			ExpectedConf: `
# Copying pipeline app-to-es to outputs
<label @APP_TO_ES>
  # Parse the logs into json
  <filter **>
    @type parser
    key_name message
    reserve_data yes
    hash_value_field structured
    <parse>
      @type json
      json_parser oj
    </parse>
  </filter>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @DEFAULT
    </store>
    <store>
      @type relabel
      @label @ES_APP_OUT
    </store>
  </match>
</label>`,
		}),
	)
})
