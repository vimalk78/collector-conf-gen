package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
)

const (
	ApplicationTags    = "kubernetes.**"
	JournalTags        = "journal.** system.var.log**"
	InfraContainerTags = "**_default_** **_kube-*_** **_openshift-*_** **_openshift_**"
	InfraTags          = InfraContainerTags + " " + JournalTags
	AuditTags          = "linux-audit.log** k8s-audit.log** openshift-audit.log** ovn-audit.log**"
)

func SourcesToInputs(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	var el []Element = make([]Element, 0)
	types := Clo.GatherSources(spec, o)
	types = Clo.AddLegacySources(types, *o)
	if types.Has(logging.InputNameApplication) {
		el = append(el, Match{
			Desc:      "Dont discard Application logs",
			MatchTags: ApplicationTags,
			MatchElement: Relabel{
				OutLabel: helpers.SourceTypeLabelName(logging.InputNameApplication),
			},
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Application logs",
			Pattern:      ApplicationTags,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el, Match{
			Desc:      "Dont discard Infrastructure logs",
			MatchTags: InfraTags,
			MatchElement: Relabel{
				OutLabel: helpers.SourceTypeLabelName(logging.InputNameInfrastructure),
			},
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Infrastructure logs",
			Pattern:      InfraTags,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	if types.Has(logging.InputNameAudit) {
		el = append(el, Match{
			Desc:      "Dont discard Audit logs",
			MatchTags: AuditTags,
			MatchElement: Relabel{
				OutLabel: helpers.SourceTypeLabelName(logging.InputNameAudit),
			},
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Audit logs",
			Pattern:      AuditTags,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	el = append(el, ConfLiteral{
		Desc:         "Send any remaining unmatched tags to stdout",
		TemplateName: "toStdout",
		Pattern:      "**",
		TemplateStr:  ToStdOut,
	})
	return el
}

var DiscardMatched string = `
{{define "discardMatched" -}}
# {{.Desc}}
<match kubernetes.**>
  @type null
</match>
{{end}}`
