package main

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

func main() {
	spec := logging.ClusterLogForwarderSpec{
		Inputs: []logging.InputSpec{},
	}
	fmt.Printf("spec: %#v\n", spec)
}
