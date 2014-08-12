// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Parse(opt *Options, c MarkdownConverter) (md Markdown) {
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
