package elements

import "text/template"

type RecordKey string
type RubyExpression string

type RecordModifier struct {
	Record     map[RecordKey]RubyExpression
	RemoveKeys []string
}

func (rm RecordModifier) Name() string {
	return "recordModifierTemplate"
}

func (rm RecordModifier) Template() string {
	return `{{define "` + rm.Name() + `"  -}}
@type record_modifier
{{if .Record -}}
<record>
{{- range $Key, $Value := .Record}}
  {{$Key}} {{$Value}}
{{- end}}
</record>
{{end -}}
{{if .RemoveKeys -}}
remove_keys {{comma_separated .RemoveKeys}}
{{end -}}
{{end}}
`
}

func (rm RecordModifier) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(rm.Template()))
}

func (rm RecordModifier) Data() interface{} {
	return rm
}
