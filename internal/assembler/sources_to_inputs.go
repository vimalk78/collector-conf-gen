package assembler

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

func (a Assembler) SourcesToInputs(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	var el []Element = make([]Element, 0)
	types := clo.GatherSources(spec)
	ApplicationTags := "kubernetes.**"
	if types.Has(logging.InputNameApplication) {
		el = append(el, Relabel{
			Desc:      "Dont discard Application logs",
			MatchTags: ApplicationTags,
			OutLabel:  sourceTypeLabelName(logging.InputNameApplication),
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Application logs",
			Pattern:      ApplicationTags,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	InfraTags := "**_default_** **_kube-*_** **_openshift-*_** **_openshift_** journal.** system.var.log**"
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el, Relabel{
			Desc:      "Dont discard Infrastructure logs",
			MatchTags: InfraTags,
			OutLabel:  sourceTypeLabelName(logging.InputNameInfrastructure),
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Infrastructure logs",
			Pattern:      InfraTags,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	AuditTags := "linux-audit.log** k8s-audit.log** openshift-audit.log**"
	if types.Has(logging.InputNameAudit) {
		el = append(el, Relabel{
			Desc:      "Dont discard Audit logs",
			MatchTags: AuditTags,
			OutLabel:  sourceTypeLabelName(logging.InputNameAudit),
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
{{define "discardMatched"}}
# {{.Desc}}
<match kubernetes.**>
  @type null
</match>{{end}}`
