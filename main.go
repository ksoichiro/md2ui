// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ksoichiro/md2ui/md"
)

func main() {
	var (
		in   = flag.String("in", "", "Input Markdown file")
		lang = flag.String("lang", "html", "Output language: Available: html")
	)
	flag.Parse()

	if *in == "" {
		fmt.Println("Input file name(-in) is required.")
		flag.Usage()
		return
	}

	var c md.MarkdownConverter
	switch *lang {
	case "android":
		c = &md.AndroidConverter{}
	case "html":
		fallthrough
	default:
		c = &md.HtmlConverter{}
	}

	result := md.ParseFile(*in, c)

	fmt.Println(strings.Join(md.Convert(&result), "\n"))
}
