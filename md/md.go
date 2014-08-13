// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

type Markdown struct {
	Elements []MarkdownElement
}

type MarkdownElement struct {
	ConverterFunc      ConverterFunc
	BlockConverterFunc BlockConverterFunc
	Values             []Inline
	Children           []MarkdownElement
	BlockQuote         bool
}

type ConverterFunc func(values []Inline) string
type BlockConverterFunc func() (openWrapper, closeWrapper string)

type Inline struct {
	Href     string
	Strong   bool
	Value    string
	Children []Inline
	NewLine  bool
	Eol      bool
}

type MarkdownConverter interface {
	ToH1(content []Inline) string
	ToH2(content []Inline) string
	ToH3(content []Inline) string
	ToH4(content []Inline) string
	ToH5(content []Inline) string
	ToH6(content []Inline) string
	ToP(content []Inline) string
	ToBlockQuote() (openWrapper, closeWrapper string)
}
