// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

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
