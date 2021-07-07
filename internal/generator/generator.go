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
	Create(*template.Template, CollectorConfType) *template.Template
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

type Generator struct {
	t CollectorConfType
}

func MakeGenerator(t CollectorConfType) *Generator {
	return &Generator{
		t: t,
	}
}

func (g *Generator) GenerateConfWithHeader(es ...Element) (string, error) {
	conf, err := g.generate(es)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{Header, conf}, "\n"), nil
}

func (g *Generator) GenerateConf(es ...Element) (string, error) {
	conf, err := g.generate(es)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(conf), nil
}

func (g *Generator) GenerateRec(t *template.Template, e Element, b *bytes.Buffer) error {
	t = e.Create(t, g.t)
	err := t.ExecuteTemplate(b, e.Name(), e.Data())
	if err != nil {
		fmt.Printf("error occured %v\n", err)
		return err
	}
	return nil
}

func (g *Generator) generate(es []Element) (string, error) {
	t := template.New("generate")
	t.Funcs(template.FuncMap{
		"generate":        g.generate,
		"compose":         g.generate,
		"indent":          indent,
		"comma_separated": comma_separated,
	})
	b := &bytes.Buffer{}
	for i, e := range es {
		if e == nil {
			e = Nil
		}
		if err := g.GenerateRec(t, e, b); err != nil {
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

//indent helper function to prefix each line of the output by N spaces
func indent(length int, in string) string {
	pad := strings.Repeat(" ", length)
	return pad + strings.Replace(in, "\n", "\n"+pad, -1)
}

func comma_separated(arr []string) string {
	return strings.Join(arr, ",")
}
