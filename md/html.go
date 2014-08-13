// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

import "strings"

type HtmlConverter struct {
}

func (c *HtmlConverter) ToH1(content []Inline) string {
	return "<h1>" + c.constructInlines(content) + "</h1>"
}

func (c *HtmlConverter) ToH2(content []Inline) string {
	return "<h2>" + c.constructInlines(content) + "</h2>"
}

func (c *HtmlConverter) ToH3(content []Inline) string {
	return "<h3>" + c.constructInlines(content) + "</h3>"
}

func (c *HtmlConverter) ToH4(content []Inline) string {
	return "<h4>" + c.constructInlines(content) + "</h4>"
}

func (c *HtmlConverter) ToH5(content []Inline) string {
	return "<h5>" + c.constructInlines(content) + "</h5>"
}

func (c *HtmlConverter) ToH6(content []Inline) string {
	return "<h6>" + c.constructInlines(content) + "</h6>"
}

func (c *HtmlConverter) ToP(content []Inline) string {
	return "<p>" + c.constructInlines(content) + "</p>"
}

func (c *HtmlConverter) ToBlockQuote() (string, string) {
	return "<blockquote>", "</blockquote>"
}

func (c *HtmlConverter) constructInlines(content []Inline) string {
	s := ""
	for _, i := range content {
		if i.NewLine {
			s += "<br />\n"
		} else if i.Strong {
			s += "<strong>" + c.constructInlines(i.Children) + "</strong>"
		} else if i.Href != "" {
			s += "<a href=\"" + i.Href + "\">" + c.constructInlines(i.Children) + "</a>"
		} else {
			s += i.Value
		}
		if i.Eol {
			s += "\n"
		}
	}
	if strings.HasSuffix(s, "\n") {
		s = strings.TrimSuffix(s, "\n")
	}
	return s
}
