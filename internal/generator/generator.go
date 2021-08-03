package generator

import (
	"bytes"
	"strings"
	"text/template"
)

type Element interface {
	Name() string
	Template() string
	Create(*template.Template) *template.Template
	Data() interface{}
}

type Section struct {
	Elements []Element
	Comment  string
}

type MultiElement interface {
	Elements() []Element
}

type InLabel = string
type OutLabel = string
type ComponentID = string

type Options map[string]string

var NoOptions = map[string]string{}

type Generator struct {
}

func MakeGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateConf(es ...Element) (string, error) {
	conf, err := g.generate(es)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(conf), nil
}

func (g *Generator) generate(es []Element) (string, error) {
	if len(es) == 0 {
		return "", nil
	}
	t := template.New("generate")
	t.Funcs(template.FuncMap{
		"compose":         func(es []Element) (string, error) { return g.generate(es) },
		"compose_one":     func(e Element) (string, error) { return g.generate([]Element{e}) },
		"indent":          indent,
		"comma_separated": commaSeparated,
	})
	b := &bytes.Buffer{}
	for i, e := range es {
		if e == nil {
			e = Nil
		}
		t = e.Create(t)
		err := t.ExecuteTemplate(b, e.Name(), e.Data())
		if err != nil {
			return "error in conf generation", err
		}
		if i < len(es)-1 {
			b.Write([]byte("\n"))
		}
	}
	return strings.TrimRight(b.String(), "\n"), nil
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
	if len(in) == 0 {
		return ""
	}
	pad := strings.Repeat(" ", length)
	inlines := strings.Split(in, "\n")
	outlines := make([]string, len(inlines))
	for i, inline := range inlines {
		// if strings.TrimSpace(inline) == "" {
		// 	outlines[i] = ""
		// } else {
		// 	outlines[i] = pad + inline
		// }
		outlines[i] = pad + inline
	}
	return strings.Join(outlines, "\n")
}

func commaSeparated(arr []string) string {
	return strings.Join(arr, ",")
}
