package source

import (
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

var HostAuditLogTemplate = `
{{define "inputSourceHostAuditTemplate" -}}
# {{.Desc}}
[sources.{{.ComponentID}}]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/audit/audit.log"]
{{end}}`

type HostAuditLog = ConfLiteral

var OpenshiftAuditLogTemplate = `
{{define "inputSourceOpenShiftAuditTemplate" -}}
# {{.Desc}}
[sources.{{.ComponentID}}]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/oauth-apiserver.audit.log"]
{{end}}
`

type OpenshiftAuditLog = ConfLiteral

var K8sAuditLogTemplate = `
{{define "inputSourceK8sAuditTemplate" -}}
# {{.Desc}}
[sources.{{.ComponentID}}]
  type = "file"
  ignore_older_secs = 600
  include = ["/var/log/kube-apiserver/audit.log"]
{{end}}
`

type K8sAuditLog = ConfLiteral
