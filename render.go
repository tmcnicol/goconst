package main

import (
	_ "embed"
	"io"
	"log"
	"strings"
	"text/template"
)

//go:embed templates/union.ts.tmpl
var tsraw string

type Data struct {
	UnionName string
	TypeName  string

	Fields []Field
}

type Field struct {
	// Documentation for the field
	Doc string
	// Name of the const
	Name string
}

type Renderer struct {
	template *template.Template
	out      io.Writer
}

func NewRenderer(out io.Writer) Renderer {
	tmpl := template.New("tsunion")
	tmpl.Funcs(template.FuncMap{
		"splitLines": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
	})
	tmpl, err := tmpl.Parse(tsraw)
	if err != nil {
		log.Fatalf("parsing %v", err)
	}

	r := Renderer{
		template: tmpl,
		out:      out,
	}

	return r
}

func (r Renderer) Render(data Data) error {
	return r.template.Execute(r.out, data)
}
