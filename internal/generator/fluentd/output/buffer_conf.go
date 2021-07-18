package output

import (
	"text/template"
)

type BufferConfig struct {
	BufferKeys           []string
	BufferPath           string
	FlushMode            string
	FlushInterval        string
	FlushThreadCount     string
	RetryType            string
	RetryWait            string
	RetryMaxInterval     string
	RetryTimeout         string
	QueuedChunkLimitSize string
	TotalLimitSize       string
	ChunkLimitSize       string
	OverflowAction       string
}

func (bc BufferConfig) Name() string {
	return "bufferConfigTemplate"
}

func (bc BufferConfig) Template() string {
	return `{{define "` + bc.Name() + `" -}}
{{- if .BufferKeys}}
<buffer {{comma_separated .BufferKeys}}>
{{- else}}
<buffer>
{{- end}}
  @type file
  path '{{.BufferPath}}'
  flush_mode {{.FlushMode}}
  flush_interval {{.FlushInterval}}
  flush_thread_count {{.FlushThreadCount}}
  flush_at_shutdown true
  retry_type {{.RetryType}}
  retry_wait {{.RetryWait}}
  retry_max_interval {{.RetryMaxInterval}}
  retry_timeout {{.RetryTimeout}}
  queued_chunks_limit_size {{.QueuedChunkLimitSize}}
  total_limit_size {{.TotalLimitSize}}
  chunk_limit_size {{.ChunkLimitSize}}
  overflow_action {{.OverflowAction}}
</buffer>
{{- end}}
`
}

func (bc BufferConfig) Data() interface{} {
	return bc
}

func (bc BufferConfig) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(bc.Template()))
}
