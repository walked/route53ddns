#!/bin/sh

GOPATH=`pwd`
export GOPATH

GOBIN=$GOPATH/bin
export GOBIN

go get 

go build 