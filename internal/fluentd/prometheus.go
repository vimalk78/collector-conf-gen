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
					Desc:      "Journal Logs go to INGRESS pipeline",
					MatchTags: "journal",
					OutLabel:  "INGRESS",
				},
				Relabel{
					Desc:      "Audit Logs go to INGRESS pipeline",
					MatchTags: "*.audit.log",
					OutLabel:  "INGRESS",
				},
				Relabel{
					Desc:      "Kubernetes Logs go to CONCAT pipeline",
					MatchTags: "kubernetes.**",
					OutLabel:  "CONCAT",
				},
			},
		},
	}
}

var EmitMetrics string = `
{{define "EmitMetrics"}}
# {{.Desc}}
<filter **>
  @type record_transformer
  enable_ruby
  <record>
	msg_size ${record.to_s.length}
  </record>
</filter>
<filter **>
  @type prometheus
  <metric>
	name cluster_logging_collector_input_record_total
	type counter
	desc The total number of incoming records
	<labels>
	  tag ${tag}
	  hostname ${hostname}
	</labels>
  </metric>
</filter>
<filter **>
  @type prometheus
  <metric>
	name cluster_logging_collector_input_record_bytes
	type counter
	desc The total bytes of incoming records
	key msg_size
	<labels>
	  tag ${tag}
	  hostname ${hostname}
	</labels>
  </metric>
</filter>
<filter **>
  @type record_transformer
  remove_keys msg_size
</filter>
{{end}}
`
