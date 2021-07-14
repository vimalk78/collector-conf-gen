package generator

type CollectorConfType1 int

const (
	CollectorConfFluentd CollectorConfType1 = iota
	CollectorConfVector
)

func SelectTemplate(e Element, ct CollectorConfType1) string {
	if ct == CollectorConfVector {
		return `{{- define "` + e.Name() + `"  -}}{{- end -}}`
	} else {
		return e.Template()
	}
}
