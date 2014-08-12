// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

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
	ToH3(content []Inline) string
	ToH4(content []Inline) string
	ToH5(content []Inline) string
	ToH6(content []Inline) string
	ToP(content []Inline) string
}
