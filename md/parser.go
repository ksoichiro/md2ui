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

type parseAttr struct {
	BlockQuoteLevel int
	Eol             bool
}

func Convert(elements *[]MarkdownElement) (result []string) {
	for _, e := range *elements {
		if 0 < len(e.Children) {
			o, c := e.BlockConverterFunc()
			result = append(result, o)
			result = append(result, Convert(&(e.Children))...)
			result = append(result, c)
		} else if e.ConverterFunc != nil {
			result = append(result, e.ConverterFunc(e.Values))
		}
	}
	return
}

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
	attr := parseAttr{BlockQuoteLevel: 0}
	index := 0
	ss := strings.Split(lines, "\n")
	parseInternal(&ss, &index, &md.Elements, &attr, c)
	return
}

func parseInternal(lines *[]string, index *int, elements *[]MarkdownElement, attr *parseAttr, c MarkdownConverter) {
	buf := []Inline{}
	for {
		if len(*lines) <= *index {
			break
		}
		s := (*lines)[*index]
		// Check current block depth
		blockQuoteLevel := 0
		for {
			if strings.HasPrefix(s, ">") {
				s = strings.TrimPrefix(s, ">")
				blockQuoteLevel++
			} else if strings.HasPrefix(s, " ") {
				s = strings.TrimPrefix(s, " ")
			} else {
				break
			}
		}
		if blockQuoteLevel > attr.BlockQuoteLevel {
			// Go down
			if 0 < len(buf) {
				*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
				buf = []Inline{}
			}
			attr.BlockQuoteLevel++
			me := MarkdownElement{BlockConverterFunc: c.ToBlockQuote, BlockQuote: true}
			// Parse as child elements
			parseInternal(lines, index, &me.Children, attr, c)
			*elements = append(*elements, me)
			continue
		} else if blockQuoteLevel < attr.BlockQuoteLevel {
			// Go up
			if strings.Trim(s, " ") == "" {
				attr.BlockQuoteLevel--
				break
			}
			*index++
		} else {
			*index++
		}
		if strings.HasPrefix(s, "######") {
			// H6
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH6, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "#####") {
			// H5
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH5, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "####") {
			// H4
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH4, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "###") {
			// H3
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH3, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "##") {
			// H2
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH2, Values: parseInline(trimHeaderChars(s))})
		} else if strings.HasPrefix(s, "#") {
			// H1
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH1, Values: parseInline(trimHeaderChars(s))})
		} else if s == "" {
			if 0 < len(buf) {
				// End of paragraph
				*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
				buf = []Inline{}
			}
		} else if 0 < len(buf) && strings.Replace(s, "=", "", -1) == "" {
			// Last line is H1
			if 1 < len(buf) {
				*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf[:len(buf)-1]})
			}
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH1, Values: buf[len(buf)-1:]})
			buf = []Inline{}
		} else if 0 < len(buf) && strings.Replace(s, "-", "", -1) == "" {
			// Last line is H2
			if 1 < len(buf) {
				*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf[:len(buf)-1]})
			}
			*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToH2, Values: buf[len(buf)-1:]})
			buf = []Inline{}
		} else {
			// P
			parseMultiline(s, &buf, attr)
		}
	}
	if 0 < len(buf) {
		*elements = append(*elements, MarkdownElement{ConverterFunc: c.ToP, Values: buf})
	}
	return
}

func parseMultiline(s string, buf *[]Inline, attr *parseAttr) {
	if strings.HasSuffix(s, "  ") {
		// New line
		*buf = append(*buf, parseInlineWithOption(strings.TrimSuffix(s, "  "), attr)...)
		*buf = append(*buf, Inline{NewLine: true})
	} else {
		attr.Eol = true
		*buf = append(*buf, parseInlineWithOption(s, attr)...)
		attr.Eol = false
	}
}

func parseInline(content string) (result []Inline) {
	return parseInlineWithOption(content, nil)
}
func parseInlineWithOption(content string, attr *parseAttr) (result []Inline) {
	tmp := content
	for {
		exp := regexp.MustCompile("^(.*)\\[([^\\]]*)\\]\\(([^\\) ]*)( +\"([^\"]*)\")?\\)(.*)$")
		groups := exp.FindStringSubmatch(tmp)
		if groups == nil || len(groups) < 1 {
			result = append(result, parseInlineStyle(tmp)...)
			break
		} else {
			result = append(result, parseInlineStyle(groups[1])...)
			result = append(result, Inline{Href: groups[3], Title: groups[5], Children: parseInlineStyle(groups[2])})
			tmp = groups[6]
		}
	}
	if attr != nil && attr.Eol && 0 < len(result) {
		result[len(result)-1].Eol = true
	}
	return result
}

func parseInlineStyle(content string) (result []Inline) {
	strongBegan := false
	buf := ""
	stars := ""
	for _, c := range strings.Split(content, "") {
		if c == "*" {
			stars += c
			if len(stars) == 2 {
				if strongBegan {
					strongBegan = false
					result = append(result, Inline{Strong: true, Children: parseInlineStyle(buf)})
					buf = ""
					stars = ""
				} else {
					strongBegan = true
				}
				stars = ""
			}
		} else {
			buf += stars + c
			stars = ""
		}
	}
	if buf != "" {
		result = append(result, Inline{Value: buf})
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
