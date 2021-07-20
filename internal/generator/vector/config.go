package vector

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

// keep no state in Conf
type Conf int

func MakeConf() Conf {
	return Conf(0)
}

func (a Conf) Assemble(spec *logging.ClusterLogForwarderSpec) []Section {
	return a.AssembleConfWithOptions(spec, &Options{})
}

func (a Conf) AssembleConfWithOptions(spec *logging.ClusterLogForwarderSpec, o *Options) []Section {
	return []Section{
		{
			a.Sources(spec, o),
			"Set of all input sources",
		},
	}
}
