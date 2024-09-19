#!/bin/bash

#部署go语言环境
wget https://golang.google.cn/dl/go1.23.1.linux-amd64.tar.gz && tar -zxf go1.23.1.linux-amd64.tar.gz -C /usr/local
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile && echo "export GOPROXY=https://goproxy.io,direct" >> /etc/profile && echo "export GOPATH=/opt/d4r" >> /etc/profile && source /etc/profile
#部署d4r
mkdir -p /opt/d4r
wget https://github.com/Xiao254182/D4R/raw/refs/heads/master/d4r.tar.gz && tar -zxf d4r.tar.gz -C /opt/d4r && cd /opt/d4r
GOOS=linux GOARCH=amd64 go build -o d4r > /dev/null 2>&1
echo "cd /opt/d4r && ./d4r" >> /usr/local/bin/d4r && chmod +x /usr/local/bin/d4r
echo "使用d4r进入系统"
