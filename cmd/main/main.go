package main

import (
	"encoding/json"
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/fluentd"
	. "github.com/vimalk78/collector-conf-gen/internal/fluentd"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PrintJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var testSpecs []struct {
	Desc string
	Spec logging.ClusterLogForwarderSpec
} = []struct {
	Desc string
	Spec logging.ClusterLogForwarderSpec
}{
	{
		Desc: "",
		Spec: logging.ClusterLogForwarderSpec{
			Inputs: []logging.InputSpec{
				{
					Name: "myapplogs",
					Application: &logging.Application{
						Namespaces: []string{"myapp"},
					},
				},
			},
			Pipelines: []logging.PipelineSpec{
				{
					InputRefs:  []string{logging.InputNameApplication, logging.InputNameInfrastructure, logging.InputNameAudit, "myapplogs"},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe",
				},
			},
		},
	},
	{
		Desc: "",
		Spec: logging.ClusterLogForwarderSpec{
			Inputs: []logging.InputSpec{
				{
					Name:  "myauditlogs",
					Audit: &logging.Audit{},
				},
				{
					Name: "myapplogs",
					Application: &logging.Application{
						Namespaces: []string{"myapp"},
						Selector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"key1": "value1",
							},
						},
					},
				},
			},
			Pipelines: []logging.PipelineSpec{
				{
					InputRefs:  []string{logging.InputNameApplication, logging.InputNameInfrastructure, logging.InputNameAudit, "myapplogs"},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe",
				},
				{
					InputRefs:  []string{"myapplogs", logging.InputNameInfrastructure},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe2",
				},
				{
					InputRefs:  []string{"myauditlogs"},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe3",
				},
			},
		},
	},
	{
		Desc: "only audit",
		Spec: logging.ClusterLogForwarderSpec{
			Pipelines: []logging.PipelineSpec{
				{
					InputRefs:  []string{logging.InputNameAudit},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe",
				},
			},
		},
	},
	{
		Desc: "only default application",
		Spec: logging.ClusterLogForwarderSpec{
			Pipelines: []logging.PipelineSpec{
				{
					InputRefs:  []string{logging.InputNameApplication},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe",
				},
			},
		},
	},
	{
		Desc: "Application with namespace",
		Spec: logging.ClusterLogForwarderSpec{
			Inputs: []logging.InputSpec{
				{
					Name: "myapplogs",
					Application: &logging.Application{
						Namespaces: []string{"myapp"},
					},
				},
			},
			Pipelines: []logging.PipelineSpec{
				{
					InputRefs:  []string{"myapplogs"},
					OutputRefs: []string{logging.OutputNameDefault},
					Name:       "my-pipe",
				},
			},
		},
	},
}

func test() {
	spec := &testSpecs[4].Spec
	spec = &logging.ClusterLogForwarderSpec{
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
				Name:       "my-pipe1",
			},
			{
				InputRefs:  []string{"myapplogs2"},
				OutputRefs: []string{logging.OutputNameDefault},
				Name:       "my-pipe2",
			},
			{
				InputRefs:  []string{logging.InputNameApplication},
				OutputRefs: []string{logging.OutputNameDefault},
				Name:       "my-pipe3",
			},
			{
				InputRefs:  []string{logging.InputNameApplication},
				OutputRefs: []string{logging.OutputNameDefault},
				Name:       "my-pipe4",
			},
		},
	}
	g := fluentd.MakeGenerator()
	s := g.MakeLoggingConf(spec)
	e := MergeSections(s)
	conf, err := fluentd.GenerateConfWithHeader(e...)
	if err != nil {
		fmt.Printf("error occured %v\n", err)
	} else {
		//	conf = ""
		fmt.Printf("%s\n", conf)
	}
	fmt.Println("--")
}

func main() {
	test()
}
