package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Element interface {
	Name() string
	Template() string
	Create(*template.Template) *template.Template
	Data() interface{}
}

type Elements []Element

type Section struct {
	Elements
	Comment string
}

type MultiElement interface {
	Elements() []Element
}

type InLabel = string
type OutLabel = string

var Header = `
## CLO GENERATED CONFIGURATION ###
# This file is a copy of the fluentd configuration entrypoint
# which should normally be supplied in a configmap.

<system>
  log_level "#{ENV['LOG_LEVEL'] || 'warn'}"
</system>

`

//indent helper function to prefix each line of the output by N spaces
func indent(length int, in string) string {
	pad := strings.Repeat(" ", length)
	return pad + strings.Replace(in, "\n", "\n"+pad, -1)
}

func comma_separated(arr []string) string {
	return strings.Join(arr, ",")
}

func GenerateConfWithHeader(es ...Element) (string, error) {
	conf, err := generate(es)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{Header, conf}, "\n"), nil
}

func GenerateConf(es ...Element) (string, error) {
	conf, err := generate(es)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(conf), nil
}

func GenerateRec(t *template.Template, e Element, b *bytes.Buffer) error {
	t = e.Create(t)
	err := t.ExecuteTemplate(b, e.Name(), e.Data())
	if err != nil {
		fmt.Printf("error occured %v\n", err)
		return err
	}
	return nil
}

func generate(es []Element) (string, error) {
	t := template.New("generate")
	t.Funcs(template.FuncMap{
		"generate": generate,
		"compose":  generate,
		"indent":   indent,
		//		"labelName":           labelName,
		//		"sourceTypelabelName": sourceTypeLabelName,
		"comma_separated": comma_separated,
	})
	b := &bytes.Buffer{}
	for i, e := range es {
		if e == nil {
			e = Nil
		}
		if err := GenerateRec(t, e, b); err != nil {
			fmt.Printf("error occured %v\n", err)
			return "", err
		}
		if i < len(es)-1 {
			b.Write([]byte("\n"))
		}
	}
	return strings.TrimSpace(b.String()), nil
}

func MergeElements(eles ...[]Element) []Element {
	merged := make([]Element, 0)
	for _, el := range eles {
		merged = append(merged, el...)
	}
	return merged
}

func MergeSections(sections []Section) []Element {
	merged := make([]Element, 0)
	for _, s := range sections {
		merged = append(merged, s.Elements...)
	}
	return merged
}
