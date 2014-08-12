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
	Level int
	Order int
	H1    string
	H2    string
	H3    string
	H4    string
	H5    string
	H6    string
	UL    string
	LI    string
	Child Markdown
	P     []string
	QUOTE []string
	CODE  []string
	HR    bool
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

	parse(&opt, c)
}

func parse(opt *Options, c MarkdownConverter) {
	filename := filepath.Join(opt.InFile)
	md, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer md.Close()

	b, _ := ioutil.ReadAll(md)
	buf := ""
	for _, s := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(s, "# ") {
			// H1
			fmt.Println(c.ToH1(c.ConvInline(strings.TrimPrefix(s, "# "))))
		} else if strings.HasPrefix(s, "## ") {
			// H2
			fmt.Println(c.ToH2(c.ConvInline(strings.TrimPrefix(s, "## "))))
		} else if s == "" {
			if buf != "" {
				// End of paragraph
				fmt.Println(c.ToP(buf))
				buf = ""
			}
		} else {
			// P
			if strings.HasSuffix(s, "  ") {
				// New line
				buf = buf + c.AddNewLine(c.ConvInline(strings.TrimSuffix(s, "  ")))
			} else {
				buf = buf + c.ConvInline(s)
			}
		}
	}
	if buf != "" {
		fmt.Println(buf)
		buf = ""
	}
}

type MarkdownConverter interface {
	AddNewLine(content string) string
	ToH1(content string) string
	ToH2(content string) string
	ToP(content string) string
	ConvInline(content string) string
}

type HtmlConverter struct {
}

func (c *HtmlConverter) AddNewLine(content string) string {
	return content + "<br />\n"
}

func (c *HtmlConverter) ToH1(content string) string {
	return "<h1>" + content + "</h1>"
}

func (c *HtmlConverter) ToH2(content string) string {
	return "<h2>" + content + "</h2>"
}

func (c *HtmlConverter) ToP(content string) string {
	return "<p>" + content + "</p>"
}

func (c *HtmlConverter) ConvInline(content string) string {
	result := content
	for {
		exp := regexp.MustCompile("^(.*)\\[([^\\]]*)\\]\\(([^\\)]*)\\)(.*)$")
		groups := exp.FindStringSubmatch(result)
		if groups == nil || len(groups) < 1 {
			return result
		} else {
			result = groups[1] + "<a href=\"" + groups[3] + "\">" + groups[2] + "</a>" + groups[4]
		}
	}
	return result
}
