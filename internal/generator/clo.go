package generator

import (
	"k8s.io/apimachinery/pkg/util/sets"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
)

/**
 These unexported methods are copied frm CLO
**/
type CLO int

const (
	IncludeLegacyForwardConfig = "includeLegacyForwardConfig"
	IncludeLegacySyslogConfig  = "includeLegacySyslogConfig"
	LegacySecureforward        = "_LEGACY_SECUREFORWARD"
	LegacySyslog               = "_LEGACY_SYSLOG"
)

//GatherSources collects the set of unique source types and namespaces
func (c CLO) GatherSources(forwarder *logging.ClusterLogForwarderSpec, op *Options) sets.String {
	types := sets.NewString()
	specs := forwarder.InputMap()
	for inputName := range logging.NewRoutes(forwarder.Pipelines).ByInput {
		if logging.ReservedInputNames.Has(inputName) {
			types.Insert(inputName) // Use name as type.
		} else if spec, ok := specs[inputName]; ok {
			if spec.Application != nil {
				types.Insert(logging.InputNameApplication)
			}
			if spec.Infrastructure != nil {
				types.Insert(logging.InputNameInfrastructure)
			}
			if spec.Audit != nil {
				types.Insert(logging.InputNameAudit)
			}
		}
	}
	return types
}

func (c CLO) AddLegacySources(types sets.String, op Options) sets.String {
	if c.IncludeLegacyForwardConfig(op) || c.IncludeLegacySyslogConfig(op) {
		types.Insert(logging.InputNameApplication)
		types.Insert(logging.InputNameInfrastructure)
		types.Insert(logging.InputNameAudit)
	}
	return types
}

func (CLO) IncludeLegacyForwardConfig(op Options) bool {
	_, ok := op[IncludeLegacyForwardConfig]
	return ok
}

func (CLO) IncludeLegacySyslogConfig(op Options) bool {
	_, ok := op[IncludeLegacySyslogConfig]
	return ok
}

var Clo CLO
