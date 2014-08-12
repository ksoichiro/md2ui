// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ParseFile(path string, c MarkdownConverter) (md Markdown) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	return Parse(string(b), c)
}

func Parse(lines string, c MarkdownConverter) (md Markdown) {
	buf := []Inline{}
	for _, s := range strings.Split(lines, "\n") {
		if strings.HasPrefix(s, "# ") {
			// H1
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH1, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "## ") {
			// H2
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH2, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "### ") {
			// H3
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH3, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "#### ") {
			// H4
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH4, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "##### ") {
			// H5
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH5, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "###### ") {
			// H6
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH6, Values: parseInline(trimHeaderChars(s))})
		} else if s == "" {
			if 0 < len(buf) {
				// End of paragraph
				md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
				buf = []Inline{}
			}
		} else if 0 < len(buf) && strings.Replace(s, "=", "", -1) == "" {
			// Last line is H1
			if 1 < len(buf) {
				md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf[:len(buf)-1]})
			}
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH1, Values: buf[len(buf)-1:]})
			buf = []Inline{}
		} else if 0 < len(buf) && strings.Replace(s, "-", "", -1) == "" {
			// Last line is H2
			if 1 < len(buf) {
				md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf[:len(buf)-1]})
			}
			md.Elements = append(md.Elements, MarkdownElement{ConverterFunc: c.ToH2, Values: buf[len(buf)-1:]})
			buf = []Inline{}
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

func trimHeaderChars(content string) string {
	result := content
	for strings.HasPrefix(result, " ") {
		result = strings.TrimPrefix(result, " ")
	}
	for strings.HasPrefix(result, "#") {
		result = strings.TrimPrefix(result, "#")
	}
	for strings.HasPrefix(result, " ") {
		result = strings.TrimPrefix(result, " ")
	}

	for strings.HasSuffix(result, " ") {
		result = strings.TrimSuffix(result, " ")
	}
	for strings.HasSuffix(result, "#") {
		result = strings.TrimSuffix(result, "#")
	}
	for strings.HasSuffix(result, " ") {
		result = strings.TrimSuffix(result, " ")
	}
	return result
}
