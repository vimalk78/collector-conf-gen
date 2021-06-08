package fluentd

import "text/template"

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
