package assembler

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type HostAuditLog struct {
	OutLabel
	Desc string
}

func (h HostAuditLog) Name() string {
	return "inputSourceHostAuditTemplate"
}

func (h HostAuditLog) Template() string {
	return `{{define "` + h.Name() + `"  -}}
# {{.Desc}}
<source>
  @type tail
  @id audit-input
  @label @{{.OutLabel}}
  path "#{ENV['AUDIT_FILE'] || '/var/log/audit/audit.log'}"
  pos_file "#{ENV['AUDIT_POS_FILE'] || '/var/log/audit/audit.log.pos'}"
  tag linux-audit.log
  <parse>
    @type viaq_host_audit
  </parse>
</source>
{{end}}`
}

func (h HostAuditLog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(h.Template()))
}

func (h HostAuditLog) Data() interface{} {
	return h
}

type OpenshiftAuditLog struct {
	OutLabel
	Desc string
}

func (o OpenshiftAuditLog) Name() string {
	return "inputSourceOpenShiftAuditTemplate"
}

func (o OpenshiftAuditLog) Template() string {
	return `{{define "` + o.Name() + `"  -}}
# {{.Desc}}
<source>
  @type tail
  @id openshift-audit-input
  @label @{{.OutLabel}}
  path /var/log/oauth-apiserver/audit.log,/var/log/openshift-apiserver/audit.log
  pos_file /var/log/oauth-apiserver.audit.log
  tag openshift-audit.log
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

func (o OpenshiftAuditLog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(o.Template()))
}

func (o OpenshiftAuditLog) Data() interface{} {
	return o
}

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
