package fluentd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

var pipeline_to_outputs = Describe("Testing Config Generation", func() {
	var f = func(g *Generator, spec logging.ClusterLogForwarderSpec) []Element {
		return MergeElements(
			g.PipelineToOutputs(&spec),
		)
	}
	DescribeTable("Pipelines(s) to Output(s)", TestGenerateConfWith(f),
		Entry("Application to single output", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "defaultoutput",
					},
				},
			},
			ExpectedConf: `
<label @DEFAULTOUTPUT>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @DEFAULT
    </store>
  </match>
</label>`,
		}),
		Entry("Application to multiple outputs", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
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
			Spec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault, "es-app-out"},
						Name:       "app-to-es",
						Labels:     map[string]string{"a": "b"},
					},
				},
			},
			ExpectedConf: `
<label @APP_TO_ES>
  # Add User Defined labels to the output record
  <filter **>
    @type record_transformer
    <record>
      openshift { "labels": {"a":"b"} }
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
			Spec: logging.ClusterLogForwarderSpec{
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
