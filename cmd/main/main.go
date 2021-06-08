package main

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/vimalk78/collector-conf-gen/internal/fluentd"
	. "github.com/vimalk78/collector-conf-gen/internal/fluentd"
)

func test() {
	s := SelectLogType(&logging.ClusterLogForwarderSpec{
		Pipelines: []logging.PipelineSpec{
			{
				InputRefs:  []string{logging.InputNameApplication, logging.InputNameInfrastructure},
				OutputRefs: []string{logging.OutputNameDefault},
				Name:       "my-pipe",
			},
		},
	})
	//Logging = append(Logging, s)
	e := MergeSections([]Section{s})
	conf, err := fluentd.GenerateConf(e...)
	if err != nil {
		fmt.Printf("error occured %v\n", err)
	} else {
		fmt.Printf("%s\n", conf)
	}
}

func main() {
	test()
}
