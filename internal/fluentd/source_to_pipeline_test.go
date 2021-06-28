package fluentd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
This test case includes only the dynamic parts of Fluentd conf. This leaves out the static parts which do not change with CLF spec.
**/
var source_to_pipline = Describe("Testing Config Generation", func() {
	var SourceToPipelines = func(g *Generator, spec logging.ClusterLogForwarderSpec) []Element {
		return MergeElements(
			g.SourceToInput(&spec),
			g.InputsToPipeline(&spec))
	}
	DescribeTable("Source(s) to Pipeline(s)", GenerateConfWith(SourceToPipelines),
		Entry("Send all log types to output by name", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs: []string{
							logging.InputNameApplication,
							logging.InputNameInfrastructure,
							logging.InputNameAudit},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Dont discard Infrastructure logs
<match **_default_** **_kube-*_** **_openshift-*_** **_openshift_** journal.** system.var.log**>
  @type relabel
  @label @_INFRASTRUCTURE
</match>

# Dont discard Audit logs
<match linux-audit.log** k8s-audit.log** openshift-audit.log**>
  @type relabel
  @label @_AUDIT
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Copying application source type to pipeline
<label @_APPLICATION>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @PIPELINE
    </store>
  </match>
</label>

# Copying infrastructure source type to pipeline
<label @_INFRASTRUCTURE>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @PIPELINE
    </store>
  </match>
</label>

# Copying audit source type to pipeline
<label @_AUDIT>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @PIPELINE
    </store>
  </match>
</label>`,
		}),
		Entry("Route Logs by Namespace(s)", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Inputs: []logging.InputSpec{
					{
						Name: "myapplogs",
						Application: &logging.Application{
							Namespaces: []string{"myapp1", "myapp2"},
						},
					},
				},
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{"myapplogs"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Discard Infrastructure logs
<match kubernetes.**>  
  @type null  
</match>

# Discard Audit logs
<match kubernetes.**>  
  @type null  
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
    <route>
      @label @PIPELINE
      <match>
        namespaces myapp1,myapp2
      </match>
    </route>
  </match>
</label>`,
		}),
		Entry("Route Logs by Labels(s)", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Inputs: []logging.InputSpec{
					{
						Name: "myapplogs",
						Application: &logging.Application{
							Selector: &v1.LabelSelector{
								MatchLabels: map[string]string{
									"key1": "value1",
									"key2": "value2",
								},
							},
						},
					},
				},
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{"myapplogs"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Discard Infrastructure logs
<match kubernetes.**>  
  @type null  
</match>

# Discard Audit logs
<match kubernetes.**>  
  @type null  
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
    <route>
      @label @PIPELINE
      <match>
        labels key1:value1,key2:value2
      </match>
    </route>
  </match>
</label>`,
		}),
		Entry("Route Logs by Namespaces(s), and Labels(s)", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Inputs: []logging.InputSpec{
					{
						Name: "myapplogs",
						Application: &logging.Application{
							Namespaces: []string{"myapp1", "myapp2"},
							Selector: &v1.LabelSelector{
								MatchLabels: map[string]string{
									"key1": "value1",
									"key2": "value2",
								},
							},
						},
					},
				},
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{"myapplogs"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Discard Infrastructure logs
<match kubernetes.**>  
  @type null  
</match>

# Discard Audit logs
<match kubernetes.**>  
  @type null  
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
    <route>
      @label @PIPELINE
      <match>
        namespaces myapp1,myapp2
        labels key1:value1,key2:value2
      </match>
    </route>
  </match>
</label>`,
		}),
		Entry("Send Logs by custom selection, and direct", ConfGenerateTest{
			Spec: logging.ClusterLogForwarderSpec{
				Inputs: []logging.InputSpec{
					{
						Name: "myapplogs",
						Application: &logging.Application{
							Namespaces: []string{"myapp1", "myapp2"},
							Selector: &v1.LabelSelector{
								MatchLabels: map[string]string{
									"key1": "value1",
									"key2": "value2",
								},
							},
						},
					},
				},
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline1",
					},
					{
						InputRefs:  []string{"myapplogs"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline2",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Discard Infrastructure logs
<match kubernetes.**>  
  @type null  
</match>

# Discard Audit logs
<match kubernetes.**>  
  @type null  
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
    <route>
      @label @PIPELINE2
      <match>
        namespaces myapp1,myapp2
        labels key1:value1,key2:value2
      </match>
    </route>
    <route>
      @label @_APPLICATION_ALL
      <match>
      </match>
    </route>
  </match>
</label>

# Copying unrouted "application" to pipelines
<label @_APPLICATION_ALL>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @PIPELINE1
    </store>
  </match>
</label>`,
		}),
		Entry("Complex Case", ConfGenerateTest{
			Desc: "Complex Case",
			Spec: logging.ClusterLogForwarderSpec{
				Inputs: []logging.InputSpec{
					{
						Name:        "myapplogs1",
						Application: &logging.Application{},
					},
					{
						Name: "myapplogs2",
						Application: &logging.Application{
							Namespaces: []string{"myapp"},
							Selector: &v1.LabelSelector{
								MatchLabels: map[string]string{
									"key1": "value1",
									"key2": "value2",
								},
							},
						},
					},
				},
				Pipelines: []logging.PipelineSpec{
					{
						InputRefs:  []string{"myapplogs1"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline1",
					},
					{
						InputRefs:  []string{"myapplogs2"},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline2",
					},
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline3",
					},
					{
						InputRefs:  []string{logging.InputNameApplication},
						OutputRefs: []string{logging.OutputNameDefault},
						Name:       "pipeline4",
					},
				},
			},
			ExpectedConf: `
# Dont discard Application logs
<match kubernetes.**>
  @type relabel
  @label @_APPLICATION
</match>

# Discard Infrastructure logs
<match kubernetes.**>  
  @type null  
</match>

# Discard Audit logs
<match kubernetes.**>  
  @type null  
</match>

# Send any remaining unmatched tags to stdout
<match **>
 @type stdout
</match>

# Routing Application to pipelines
<label @_APPLICATION>
  <match **>
    @type label_router
    <route>
      @label @PIPELINE2
      <match>
        namespaces myapp
        labels key1:value1,key2:value2
      </match>
    </route>
    <route>
      @label @_APPLICATION_ALL
      <match>
      </match>
    </route>
  </match>
</label>

# Copying unrouted "application" to pipelines
<label @_APPLICATION_ALL>
  <match **>
    @type copy
    <store>
      @type relabel
      @label @PIPELINE1
    </store>
    <store>
      @type relabel
      @label @PIPELINE3
    </store>
    <store>
      @type relabel
      @label @PIPELINE4
    </store>
  </match>
</label>`,
		}),
	)
})
