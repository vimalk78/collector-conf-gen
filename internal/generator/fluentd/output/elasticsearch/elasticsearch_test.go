package elasticsearch

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

func TestFluendConfGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fluend Conf Generation")
}

var elasticsearch_store_test = PDescribe("", func() {
	var f = func(clspec logging.ClusterLoggingSpec, clfspec logging.ClusterLogForwarderSpec) []Element {
		return Store(clfspec.Outputs[0], &Options{})
	}
	DescribeTable("Elassticsearch store", TestGenerateConfWith(f),
		Entry("", ConfGenerateTest{
			CLFSpec: logging.ClusterLogForwarderSpec{
				Outputs: []logging.OutputSpec{
					{
						Type: logging.OutputTypeElasticsearch,
						Name: "es-1",
					},
				},
			},
			ExpectedConf: ``,
		}))
})
