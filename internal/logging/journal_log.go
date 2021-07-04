package logging

import (
	"text/template"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

type JournalLog struct {
	OutLabel
	Desc string
}

func (j JournalLog) Name() string {
	return "inputSourceJournalTemplate"
}

func (j JournalLog) Template() string {
	return `{{define "` + j.Name() + `"  -}}
# {{.Desc}}
<source>
  @type systemd
  @id systemd-input
  @label @{{.OutLabel}}
  path '/var/log/journal'
  <storage>
    @type local
    persistent true
    # NOTE: if this does not end in .json, fluentd will think it
    # is the name of a directory - see fluentd storage_local.rb
    path '/var/log/journal_pos.json'
  </storage>
  matches "#{ENV['JOURNAL_FILTERS_JSON'] || '[]'}"
  tag journal
  read_from_head "#{if (val = ENV.fetch('JOURNAL_READ_FROM_HEAD','')) && (val.length > 0); val; else 'false'; end}"
</source>
{{end}}`
}

func (j JournalLog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(j.Template()))
}

func (j JournalLog) Data() interface{} {
	return j
}
