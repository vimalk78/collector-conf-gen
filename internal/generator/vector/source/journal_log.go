package source

import (
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

var JournalLogTemplate = `
{{define "inputSourceJournalTemplate" -}}
[sources.{{.ComponentID}}]
  type = "journald"
{{end}}`

type JournalLog = ConfLiteral
