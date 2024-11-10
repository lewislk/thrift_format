#!/bin/bash
RUN_NAME="thrift_format"

# windows
go build -o output/bin/${RUN_NAME}.exe

# linux/mac os
#go build -o output/bin/${RUN_NAME}