#!/usr/bin/env bash

PROG=_md2ui

go build -o $PROG

./$PROG -in testdata/test.md -lang android > testdata/android/wrapper/src/main/res/layout/activity_top.xml

if [ $? -ne 0 ]; then
    echo "Failed to run"
    exit 1
fi

pushd testdata/android/wrapper/ > /dev/null
./gradlew installDebug
popd > /dev/null

rm -f $PROG

