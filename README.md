# md2ui

Parse Markdown and generate some kind of UI.

This is still an experimental project.

## Install

```sh
$ go get github.com/ksoichiro/md2ui
```

## Usage

```sh
$ md2ui -in testdata/test.md
```

This will print HTML from `test.md`.

## Converter - what kind of UI can we use?

Converter can be specify with option `-lang`.

* HTML(`-lang html`)
* Android XML (`-lang android`)

## License

Copyright (c) 2014 Soichiro Kashima  
Licensed under MIT license.  
See the bundled [LICENSE](https://github.com/ksoichiro/rdotm/blob/master/LICENSE) file for details.
