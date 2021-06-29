package fluentd

import logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"

func (g *Generator) SourceToInput(spec *logging.ClusterLogForwarderSpec) []Element {
	var el []Element = make([]Element, 0)
	types := clo.GatherSources(spec)
	ApplicationPattern := "kubernetes.**"
	if types.Has(logging.InputNameApplication) {
		el = append(el, Relabel{
			Desc:     "Dont discard Application logs",
			Pattern:  ApplicationPattern,
			OutLabel: "_APPLICATION",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Application logs",
			Pattern:      ApplicationPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	InfraPattern := "**_default_** **_kube-*_** **_openshift-*_** **_openshift_** journal.** system.var.log**"
	if types.Has(logging.InputNameInfrastructure) {
		el = append(el, Relabel{
			Desc:     "Dont discard Infrastructure logs",
			Pattern:  InfraPattern,
			OutLabel: "_INFRASTRUCTURE",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Infrastructure logs",
			Pattern:      InfraPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	AuditPattern := "linux-audit.log** k8s-audit.log** openshift-audit.log**"
	if types.Has(logging.InputNameAudit) {
		el = append(el, Relabel{
			Desc:     "Dont discard Audit logs",
			Pattern:  AuditPattern,
			OutLabel: "_AUDIT",
		})
	} else {
		el = append(el, ConfLiteral{
			Desc:         "Discard Audit logs",
			Pattern:      AuditPattern,
			TemplateName: "discardMatched",
			TemplateStr:  DiscardMatched,
		})
	}
	el = append(el, ConfLiteral{
		Desc:         "Send any remaining unmatched tags to stdout",
		TemplateName: "toStdout",
		TemplateStr: `
{{define "toStdout" -}}
# {{.Desc}}
<match **>
 @type stdout
</match>
{{end -}}`,
	})
	return el
}

var DiscardMatched string = `
{{define "discardMatched" -}}
# {{.Desc}}
<match kubernetes.**>
  @type null
</match>
{{end}}
`
