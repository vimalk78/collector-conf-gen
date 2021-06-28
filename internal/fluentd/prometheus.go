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

func (g *Generator) PrometheusMetrics() []Element {
	return []Element{
		Pipeline{
			InLabel: "MEASURE",
			Desc:    "Increment Prometheus metrics",
			SubElements: []Element{
				ConfLiteral{
					TemplateName: "EmitMetrics",
					Desc:         "xx",
					TemplateStr:  EmitMetrics,
				},
				Relabel{
					Desc:     "Journal Logs go to INGRESS pipeline",
					Pattern:  "journal",
					OutLabel: "INGRESS",
				},
				Relabel{
					Desc:     "Audit Logs go to INGRESS pipeline",
					Pattern:  "*.audit.log",
					OutLabel: "INGRESS",
				},
				Relabel{
					Desc:     "Kubernetes Logs go to CONCAT pipeline",
					Pattern:  "kubernetes.**",
					OutLabel: "CONCAT",
				},
			},
		},
	}
}
