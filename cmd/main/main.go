package main

import (
	"encoding/json"
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/generator"
	loggen "github.com/vimalk78/collector-conf-gen/internal/logging"
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
	conf, _ := generator.GenerateConf(
		generator.MergeSections(
			loggen.MakeGenerator().MakeLoggingConf(&spec))...)
	fmt.Printf("conf:\n%s\n", conf)
}

func main() {
	test()
}
