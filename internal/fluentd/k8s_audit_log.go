package fluentd

import "text/template"

type K8sAuditLog struct {
	OutLabel
	Desc string
}

func (k K8sAuditLog) Name() string {
	return "inputSourceK8sAuditTemplate"
}

func (k K8sAuditLog) Template() string {
	return `{{define "` + k.Name() + `"  -}}
# {{.Desc}}
<source>
  @type tail
  @id k8s-audit-input
  @label @{{.OutLabel}}
  path "#{ENV['K8S_AUDIT_FILE'] || '/var/log/kube-apiserver/audit.log'}"
  pos_file "#{ENV['K8S_AUDIT_POS_FILE'] || '/var/log/kube-apiserver/audit.log.pos'}"
  tag k8s-audit.log
  <parse>
    @type json
    time_key requestReceivedTimestamp
    # In case folks want to parse based on the requestReceivedTimestamp key
    keep_time_key true
    time_format %Y-%m-%dT%H:%M:%S.%N%z
  </parse>
</source>
{{end}}`
}

func (k K8sAuditLog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(k.Template()))
}

func (k K8sAuditLog) Data() interface{} {
	return k
}
