// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package md

import "strings"

type AndroidConverter struct {
}

type textViewAttr struct {
	Style string
}

const (
	styleH1 = "MarkdownH1"
	styleH2 = "MarkdownH2"
	styleH3 = "MarkdownH3"
	styleH4 = "MarkdownH4"
	styleH5 = "MarkdownH5"
	styleH6 = "MarkdownH6"
	styleP  = "MarkdownP"
)

func (c *AndroidConverter) Header() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<ScrollView xmlns:android="http://schemas.android.com/apk/res/android"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    >
<LinearLayout
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:orientation="vertical"
    >`
}

func (c *AndroidConverter) Footer() string {
	return `</LinearLayout>
</ScrollView>`
}

func (c *AndroidConverter) ToH1(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH1})
}

func (c *AndroidConverter) ToH2(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH2})
}

func (c *AndroidConverter) ToH3(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH3})
}

func (c *AndroidConverter) ToH4(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH4})
}

func (c *AndroidConverter) ToH5(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH5})
}

func (c *AndroidConverter) ToH6(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleH6})
}

func (c *AndroidConverter) ToP(content []Inline) string {
	return c.toTextView(content, textViewAttr{Style: styleP})
}

func (c *AndroidConverter) ToBlockQuote() (openTag string, closeTag string) {
	openTag = "<LinearLayout\n"
	openTag += "    android:layout_width=\"match_parent\"\n"
	openTag += "    android:layout_height=\"wrap_content\"\n"
	openTag += "    android:orientation=\"vertical\"\n"
	openTag += "    android:padding=\"8dp\"\n"
	openTag += "    >"
	closeTag = "</LinearLayout>"
	return
}

func (c *AndroidConverter) constructInlines(content []Inline) string {
	s := ""
	for _, i := range content {
		if 0 < len(i.Children) {
			s += c.constructInlines(i.Children)
		} else if i.NewLine {
			s += "\\n"
		} else {
			s += i.Value
		}
	}
	if strings.HasSuffix(s, "\n") {
		s = strings.TrimSuffix(s, "\n")
	}
	return s
}

func (c *AndroidConverter) toTextView(content []Inline, attr textViewAttr) string {
	result := "<TextView\n"
	result += "    android:layout_width=\"wrap_content\"\n"
	result += "    android:layout_height=\"wrap_content\"\n"
	result += "    android:text=\"" + c.constructInlines(content) + "\"\n"
	if attr.Style != "" {
		result += "    style=\"@style/" + attr.Style + "\"\n"
	}
	result += "    />\n"
	return result
}
