// Program generates ../2025-07-28-modernize-go.md article from 2025-07-28-modernize-go.tmpl.
// Run `go generate ./...` in this directory to update contents.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go run $GOFILE

type Go struct {
	Version  string
	Sections []Section
}

type Section struct {
	Name string

	Header   string
	Body     string
	Benefit  string
	Command  template.HTML
	Before   template.HTML
	After    template.HTML
	Examples string
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "If print debug statements")
	flag.Parse()

	const articleFilename = "2025-07-28-modernize-go"
	tmplFilename := articleFilename + ".tmpl"
	tmpl, err := template.New(tmplFilename).ParseFiles(tmplFilename)
	if err != nil {
		log.Fatal(err)
	}

	gos := []Go{
		{
			Version: "1.24",
		},
		{
			Version: "1.23",
		},
		{
			Version: "1.22",
			Sections: []Section{
				{
					Name:    "forloop",
					Header:  "Remove redundant loop variables",
					Body:    "A construct like `tc := tc` in `for` loops is not needed anymore and we can remove it.\nSee [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview) for details.",
					Benefit: "Saves one line of code.",
					Examples: `- [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser/pull/4856/files#diff-3756619488c8c0f0c0300fc0cdcfecbb39c2a7bcb4fe4b3ac5305c6057512986L486)
- [kubernetes-sigs/kueue](https://github.com/kubernetes-sigs/kueue/pull/1946/files#diff-22ad2263a86a607fd28df7741c704614d0f34e208b5270153aa39427e4325fb3L203)
- [IBM/sarama](https://github.com/IBM/sarama/pull/3214/files#diff-cb488ad8239edeaaf8b0c1f469cc15c03fde53cbf22ee996e2f3922b3cc6a0c9L426)
- [google/go-github](https://github.com/google/go-github/pull/3537/files#diff-0f446fb8e4e16b655368f9f1c774d667d5528c9b3103f35481f704e2e33a925fL292)
- [go-critic/go-critic](https://github.com/go-critic/go-critic/pull/1459/files#diff-c2dfb8c940e1232344ce37c2a5942712765d9acf23d43c89345feb81fdbeeb13L43)
- [99designs/gqlgen](https://github.com/99designs/gqlgen/pull/3387/files#diff-fa4826c514673a47321901386ae757f00b2faa73d1433d8dacfc836f4928829aL44)
- [air-verse/air](https://github.com/air-verse/air/pull/682/files#diff-0c22297be1ae696feec687c4dc3d1f425a6ff6c7dfd47d1d2a2275c32d3da14aL96)`,
				},
				{
					Name:    "forrange",
					Header:  "Simplify `for` range loops",
					Body:    "\"For\" loops may now range over integers.\nSee [For-range over integers in Go 1.22](https://go.dev/ref/spec#For_range) for details.",
					Benefit: "Improves readability and less symbols to type.",
				},
			},
		},
		{
			Version: "1.21",
		},
		{
			Version: "1.20",
		},
		{
			Version: "1.19",
		},
		{
			Version: "1.18",
		},
		{
			Version: "1.17",
		},
		{
			Version: "1.16",
		},
		{
			Version: "1.15",
		},
		{
			Version: "1.14",
		},
		{
			Version: "1.13",
		},
		{
			Version: "1.12",
		},
		{
			Version: "1.11",
		},
		{
			Version: "1.10",
		},
		{
			Version: "1.9",
		},
		{
			Version: "1.8",
		},
		{
			Version: "1.7",
		},
		{
			Version: "1.6",
		},
		{
			Version: "1.5",
		},
		{
			Version: "1.4",
		},
		{
			Version: "1.3",
		},
		{
			Version: "1.2",
		},
		{
			Version: "1.1",
		},
	}

	for _, goElem := range gos {
		for i := range goElem.Sections {
			section := goElem.Sections[i]
			commandFilename := filepath.Join(goElem.Version, commandFilename(section.Name))
			command, err := extractCommandContent(commandFilename)
			if err != nil {
				log.Fatal(err)
			}
			goElem.Sections[i].Command = command

			beforeFilename := filepath.Join(goElem.Version, beforeFilename(section.Name))
			before, err := extractGoContent(beforeFilename)
			if err != nil {
				log.Fatal(err)
			}
			goElem.Sections[i].Before = before

			afterFilename := filepath.Join(goElem.Version, afterFilename(section.Name))
			after, err := extractGoContent(afterFilename)
			if err != nil {
				log.Fatal(err)
			}
			goElem.Sections[i].After = template.HTML(after)
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, gos)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println(buf.String())
	}

	// workaround
	b := bytes.Replace(buf.Bytes(), []byte("&lt;!--more--&gt;"), []byte("<!--more-->"), 1)

	err = os.WriteFile(filepath.Join("..", articleFilename+".md"), b, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func commandFilename(name string) string {
	return name + "_command.sh"
}

func beforeFilename(name string) string {
	return name + "_before.go"
}

func afterFilename(name string) string {
	return name + "_after.go"
}

func extractCommandContent(filename string) (template.HTML, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	command, ok := strings.CutPrefix(string(contents), "#! /bin/sh")
	if !ok {
		return "", errors.New("missed sheband")
	}
	command = strings.TrimSpace(command)
	return template.HTML(command), nil
}

func extractGoContent(filename string) (template.HTML, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	const begin = "// << begin >>\n"
	_, after, ok := strings.Cut(string(contents), begin)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", begin, filename)
	}
	const end = "\n\t// << end >>"
	before, _, ok := strings.Cut(after, end)
	if !ok {
		return "", fmt.Errorf("missed %q in file %q", end, filename)
	}
	return template.HTML(before), nil
}
