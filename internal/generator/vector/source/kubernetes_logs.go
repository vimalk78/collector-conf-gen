package source

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type KubernetesLogs struct {
	ComponentID
	Desc         string
	ExcludePaths string
}

func (kl KubernetesLogs) Name() string {
	return "k8s_logs_template"
}

func (kl KubernetesLogs) Template() string {
	return `{{define "` + kl.Name() + `" -}}
# {{.Desc}}
[sources.{{.ComponentID}}]
  auto_partial_merge = true
  exclude_paths_glob_patterns = {{.ExcludePaths}}
{{end}}`
}

func (kl KubernetesLogs) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(kl.Template()))
}

func (kl KubernetesLogs) Data() interface{} {
	return kl
}
