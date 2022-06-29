#!/bin/sh

mkdir "releases"

# 可选参数-ldflags 是编译选项：
#   -s -w 去掉调试信息，可以减小构建后文件体积。

# 【darwin/amd64】
echo "start build darwin/amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w --extldflags "-static -fpic"' -o ./releases/showdocdb-darwin-amd64
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/showdocdb-darwin-amd64-cgo

# 【linux/amd64】
echo "start build linux/amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w --extldflags "-static -fpic"' -o ./releases/showdocdb-linux-amd64
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/showdocdb-linux-amd64-cgo

# 【windows/amd64】
echo "start build windows/amd64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w --extldflags "-static -fpic"' -o ./releases/showdocdb-windows-amd64.exe
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/showdocdb-windows-amd64-cgo.exe

echo "Congratulations,all build success!!!"