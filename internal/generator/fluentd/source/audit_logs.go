package source

import (
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

var HostAuditLogTemplate = `
{{define "inputSourceHostAuditTemplate" -}}
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

type HostAuditLog = ConfLiteral

var OpenshiftAuditLogTemplate = `
{{define "inputSourceOpenShiftAuditTemplate" -}}
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
{{end}}
`

type OpenshiftAuditLog = ConfLiteral

var K8sAuditLogTemplate = `
{{define "inputSourceK8sAuditTemplate" -}}
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
{{end}}
`

type K8sAuditLog = ConfLiteral
