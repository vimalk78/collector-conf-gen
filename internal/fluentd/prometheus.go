package fluentd

import "text/template"

type PrometheusMonitor struct {
	Desc string
}

func (p PrometheusMonitor) Name() string {
	return "inputSourceHostAuditTemplate"
}

func (p PrometheusMonitor) Template() string {
	return `{{define "` + p.Name() + `"  -}}
# {{.Desc}}
<source>
  @type prometheus
  bind "#{ENV['POD_IP']}"
  <ssl>
    enable true
    certificate_path "#{ENV['METRICS_CERT'] || '/etc/fluent/metrics/tls.crt'}"
    private_key_path "#{ENV['METRICS_KEY'] || '/etc/fluent/metrics/tls.key'}"
  </ssl>
</source>

<source>
  @type prometheus_monitor
  <labels>
    hostname ${hostname}
  </labels>
</source>

# excluding prometheus_tail_monitor
# since it leaks namespace/pod info
# via file paths

# This is considered experimental by the repo
<source>
  @type prometheus_output_monitor
  <labels>
    hostname ${hostname}
  </labels>
</source>
{{end}}`
}

func (p PrometheusMonitor) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(p.Template()))
}

func (p PrometheusMonitor) Data() interface{} {
	return p
}
