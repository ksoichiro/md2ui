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
		in  = flag.String("in", "", "Input Markdown file")
		out = flag.String("out", "out", "Output directory for generated codes")
	)
	flag.Parse()

	opt := Options{
		InFile: *in,
		OutDir: *out,
	}

	if *in == "" {
		fmt.Println("Input file name(-in) is required.")
		return
	}

	parse(&opt)
}

func parse(opt *Options) {
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
		// FIXME This is just test codes. Not layout
		if strings.HasPrefix(s, "# ") {
			// H1
			fmt.Println(toH1(convInline(strings.TrimPrefix(s, "# "))))
		} else if strings.HasPrefix(s, "## ") {
			// H2
			fmt.Println(toH2(convInline(strings.TrimPrefix(s, "## "))))
		} else if s == "" {
			if buf != "" {
				// End of paragraph
				fmt.Println(toP(buf))
				buf = ""
			}
		} else {
			// P
			if strings.HasSuffix(s, "  ") {
				// New line
				buf = buf + addNewLine(convInline(strings.TrimSuffix(s, "  ")))
			} else {
				buf = buf + convInline(s)
			}
		}
	}
	if buf != "" {
		fmt.Println(buf)
		buf = ""
	}
}

func addNewLine(content string) string {
	return content + "<br />\n"
}

func toH1(content string) string {
	return "<h1>" + content + "</h1>"
}

func toH2(content string) string {
	return "<h2>" + content + "</h2>"
}

func toP(content string) string {
	return "<p>" + content + "</p>"
}

func convInline(content string) string {
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
