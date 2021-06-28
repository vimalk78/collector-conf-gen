package main

import (
	"encoding/json"
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/fluentd"
)

func PrintJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func test() {
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
	g := fluentd.MakeGenerator()
	s := g.MakeLoggingConf(&spec)
	e := fluentd.MergeSections(s)
	conf, _ := fluentd.GenerateConf(e...)
	fmt.Printf("conf:\n%s\n", conf)
}

func main() {
	test()
}
