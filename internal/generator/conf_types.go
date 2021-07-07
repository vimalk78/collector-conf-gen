package generator

type CollectorConfType int

const (
	CollectorConfFluentd CollectorConfType = iota
	CollectorConfVector
)

func SelectTemplate(e Element, ct CollectorConfType) string {
	if ct == CollectorConfVector {
		return `{{- define "` + e.Name() + `"  -}}{{- end -}}`
	} else {
		return e.Template()
	}
}
