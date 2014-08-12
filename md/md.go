// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

type Options struct {
	InFile string
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

type MarkdownConverter interface {
	ToH1(content []Inline) string
	ToH2(content []Inline) string
	ToP(content []Inline) string
}
