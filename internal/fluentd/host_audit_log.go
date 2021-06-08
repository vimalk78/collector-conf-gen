package fluentd

import "text/template"

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
