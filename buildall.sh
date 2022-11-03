#!/bin/bash -x

go test

go build

GOOS=windows go build
