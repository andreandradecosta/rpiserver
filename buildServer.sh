#!/bin/bash

go build -ldflags "-X main.buildInfo=`git rev-parse --short HEAD`" -o bin/rpiserver rpiserver.go
