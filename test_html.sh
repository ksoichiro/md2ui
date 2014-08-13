#!/usr/bin/env bash

PROG=_md2ui

go build -o $PROG

./$PROG -in testdata/test.md

rm -f $PROG

