package main

import (
	"encoding/json"
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd"
)

func PrintJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func testFluentd() {
	spec := logging.ClusterLogForwarderSpec{
		Pipelines: []logging.PipelineSpec{
			{
				InputRefs: []string{
					logging.InputNameApplication,
				},
				OutputRefs: []string{logging.OutputNameDefault},
				Name:       "pipeline",
			},
		},
	}
	g := generator.MakeGenerator()
	conf, _ := g.GenerateConfWithHeader(
		generator.MergeSections(
			fluentd.MakeConf().Assemble(&spec))...)
	fmt.Printf("conf:\n%s\n", conf)
}

func testVector() {
	spec := logging.ClusterLogForwarderSpec{
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
	}
	_ = spec
}

func main() {
	testFluentd()
	//testVector()
}
