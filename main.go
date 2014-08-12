// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package main

import (
	"flag"
	"fmt"

	"github.com/ksoichiro/md2ui/md"
)

func main() {
	var (
		in   = flag.String("in", "", "Input Markdown file")
		lang = flag.String("lang", "html", "Output language: Available: html")
	)
	flag.Parse()

	opt := md.Options{
		InFile: *in,
		Lang:   *lang,
	}

	if *in == "" {
		fmt.Println("Input file name(-in) is required.")
		flag.Usage()
		return
	}

	var c md.MarkdownConverter
	switch opt.Lang {
	case "html":
		fallthrough
	default:
		c = &md.HtmlConverter{}
	}

	result := md.Parse(&opt, c)

	for _, e := range result.Elements {
		fmt.Println(e.ConverterFunc(e.Values))
	}
}
