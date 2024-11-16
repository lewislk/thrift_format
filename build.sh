#!/bin/bash
RUN_NAME="thrift_format"

# go通过使用不同的环境变量可以打包不同平台运行的程序

# mac下的环境变量
# intel cpu
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/bin/${RUN_NAME}-mac-amd64
# apple cpu
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/bin/${RUN_NAME}-mac-arm64

# linux的环境变量
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/bin/${RUN_NAME}-linux-amd64

# windows的环境变量
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o output/bin/${RUN_NAME}.exe