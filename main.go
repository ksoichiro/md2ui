// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type Options struct {
	InFile string
	OutDir string
	Lang   string
}

type Markdown struct {
	Elements []MarkdownElement
}

type MarkdownElement struct {
	ConverterFunc ConverterFunc
	Values        []Inline
}

type ConverterFunc func(values []Inline) string

type Inline struct {
	Href    string
	Value   string
	NewLine bool
}

func main() {
	var (
		in   = flag.String("in", "", "Input Markdown file")
		out  = flag.String("out", "out", "Output directory for generated codes")
		lang = flag.String("lang", "html", "Output language: Available: html")
	)
	flag.Parse()

	opt := Options{
		InFile: *in,
		OutDir: *out,
		Lang:   *lang,
	}

	if *in == "" {
		fmt.Println("Input file name(-in) is required.")
		return
	}

	var c MarkdownConverter
	switch opt.Lang {
	case "html":
		fallthrough
	default:
		c = &HtmlConverter{}
	}

	md := parse(&opt, c)

	for _, e := range md.Elements {
		fmt.Println(e.ConverterFunc(e.Values))
	}
}

func parse(opt *Options, c MarkdownConverter) (md Markdown) {
	filename := filepath.Join(opt.InFile)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)
	buf := []Inline{}
	for _, s := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(s, "# ") {
			// H1
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH1, Values: parseInline(strings.TrimPrefix(s, "# "))})
		} else if strings.HasPrefix(s, "## ") {
			// H2
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH2, Values: parseInline(strings.TrimPrefix(s, "## "))})
		} else if s == "" {
			if 0 < len(buf) {
				// End of paragraph
				md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
				buf = []Inline{}
			}
		} else {
			// P
			if strings.HasSuffix(s, "  ") {
				// New line
				buf = append(buf, parseInline(strings.TrimSuffix(s, "  "))...)
				buf = append(buf, Inline{NewLine: true})
			} else {
				buf = append(buf, parseInline(s)...)
			}
		}
	}
	if 0 < len(buf) {
		md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
	}
	return
}

func parseInline(content string) (result []Inline) {
	tmp := content
	for {
		exp := regexp.MustCompile("^(.*)\\[([^\\]]*)\\]\\(([^\\)]*)\\)(.*)$")
		groups := exp.FindStringSubmatch(tmp)
		if groups == nil || len(groups) < 1 {
			return append(result, Inline{Value: tmp})
		} else {
			result = append(result, Inline{Value: groups[1]})
			result = append(result, Inline{Href: groups[3], Value: groups[2]})
			tmp = groups[4]
		}
	}
	return result
}

type MarkdownConverter interface {
	ToH1(content []Inline) string
	ToH2(content []Inline) string
	ToP(content []Inline) string
}

type HtmlConverter struct {
}

func (c *HtmlConverter) ToH1(content []Inline) string {
	return "<h1>" + c.constructInlines(content) + "</h1>"
}

func (c *HtmlConverter) ToH2(content []Inline) string {
	return "<h2>" + c.constructInlines(content) + "</h2>"
}

func (c *HtmlConverter) ToP(content []Inline) string {
	return "<p>" + c.constructInlines(content) + "</p>"
}

func (c *HtmlConverter) constructInlines(content []Inline) string {
	s := ""
	for _, i := range content {
		if i.NewLine {
			s += "<br />"
		} else if i.Href != "" {
			s += "<a href=\"" + i.Href + "\">" + i.Value + "</a>"
		} else {
			s += i.Value
		}
	}
	return s
}
