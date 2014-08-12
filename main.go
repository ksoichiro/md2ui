// Copyright (c) 2014 Soichiro Kashima
// Licensed under MIT license.

package main

import (
	"flag"
	"fmt"

	"github.com/ksoichiro/md2ui/converter"
)

const (
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type Options struct {
	InFile string
	OutDir string
	Lang   string
}

func main() {
	var (
		in   = flag.String("in", "", "Input Markdown file")
		out  = flag.String("out", "out", "Output directory for generated codes")
		lang = flag.String("lang", "html", "Output language: Available: html")
	)
	flag.Parse()

	opt := Options{
		InFile: *in,
		OutDir: *out,
		Lang:   *lang,
	}

	if *in == "" {
		fmt.Println("Input file name(-in) is required.")
		return
	}

	var c converter.MarkdownConverter
	switch opt.Lang {
	case "html":
		fallthrough
	default:
		c = &converter.HtmlConverter{}
	}

	md := parse(&opt, c)

	for _, e := range md.Elements {
		fmt.Println(e.ConverterFunc(e.Values))
	}
}
