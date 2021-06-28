package fluentd

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

type ConfGenerateTest struct {
	Desc         string
	Spec         logging.ClusterLogForwarderSpec
	ExpectedConf string
}

type GenerateFunc func(*Generator, logging.ClusterLogForwarderSpec) []Element

func GenerateConfWith(gf GenerateFunc) func(ConfGenerateTest) {
	return func(testcase ConfGenerateTest) {
		g := MakeGenerator()
		e := gf(g, testcase.Spec)
		conf, err := GenerateConf(e...)
		Expect(err).To(BeNil())
		diff := cmp.Diff(
			strings.Split(strings.TrimSpace(testcase.ExpectedConf), "\n"),
			strings.Split(strings.TrimSpace(conf), "\n"))
		if diff != "" {
			fmt.Println(conf)
			fmt.Printf("diff: %s", diff)
		}
		Expect(diff).To(Equal(""))
	}
}
