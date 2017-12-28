#!/usr/bin/env bash

case $1 in
  pi)
    echo "Building for raspberry pi - arm"
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o webos *.go
    ;;
  *)
    go build -o webos *.go
esac
