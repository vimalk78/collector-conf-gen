package syslog

import (
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	//. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
)

type Syslog struct {
	Desc    string
	StoreID string
}

func (s Syslog) Name() string {
	return "syslogTemplate"
}

func (s Syslog) Template() string {
	return `{{define "` + s.Name() + `" -}}
@type remote_syslog
@id {{.StoreID}}
host {{.Host}}
port {{.Port}}
rfc {{.Rfc}}
facility {{.Facility}}
severity {{.Severity}}

{{if .HasAppName -}}
appname {{.AppName}}
{{end -}}

{{if .Target.Syslog.MsgID -}}
msgid {{.MsgID}}
{{end -}}

{{if .Target.Syslog.ProcID -}}
procid {{.ProcID}}
{{end -}}

{{if .HasTag -}}
program {{.Tag}}
{{end -}}

protocol {{.Protocol}}
packet_size 4096
hostname "#{ENV['NODE_NAME']}"

{{ if .Target.Secret -}}
tls true
ca_file '{{ .SecretPath "ca-bundle.crt"}}'
verify_mode true
{{ end -}}

{{ if (eq .Protocol "tcp") -}}
timeout 60
timeout_exception true
keep_alive true
keep_alive_idle 75
keep_alive_cnt 9
keep_alive_intvl 7200
{{ end -}}

{{if .PayloadKey -}}
<format>
@type single_value
message_key {{.PayloadKey}}
</format>
{{end -}}
<buffer {{.ChunkKeys}}>
@type file
path '{{.BufferPath}}'
flush_mode {{.FlushMode}}
flush_interval {{.FlushInterval}}
flush_thread_count {{.FlushThreadCount}}
flush_at_shutdown true
retry_type {{.RetryType}}
retry_wait {{.RetryWait}}
retry_max_interval {{.RetryMaxInterval}}
{{.RetryTimeout}}
queued_chunks_limit_size "#{ENV['BUFFER_QUEUE_LIMIT'] || '32' }"
{{- if .TotalLimitSize }}
total_limit_size {{.TotalLimitSize}}
{{- else }}
total_limit_size "#{ENV['TOTAL_LIMIT_SIZE'] ||  8589934592 }" #8G
{{- end }}
{{- if .ChunkLimitSize }}
chunk_limit_size {{.ChunkLimitSize}}
{{- else }}
chunk_limit_size "#{ENV['BUFFER_SIZE_LIMIT'] || '8m'}"
{{- end }}
overflow_action {{.OverflowAction}}
</buffer>
</store>
{{end}}
`
}

func (s Syslog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(s.Template()))
}

func (s Syslog) Data() interface{} {
	return s
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	return nil
}
