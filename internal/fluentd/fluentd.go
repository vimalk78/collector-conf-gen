package fluentd

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

type MultiElement interface {
	Elements() []Element
}

type InLabel string
type OutLabel string

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

func GenerateConf(es ...Element) (string, error) {
	conf, err := generate(es)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{Header, conf}, "\n"), nil
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
		"indent":   indent,
	})
	b := &bytes.Buffer{}
	for i, e := range es {
		if err := GenerateRec(t, e, b); err != nil {
			fmt.Printf("error occured %v\n", err)
			return "", err
		}
		if i < len(es)-1 {
			// templates must add their own new line
			b.Write([]byte("\n"))
		}
	}
	return b.String(), nil
}
