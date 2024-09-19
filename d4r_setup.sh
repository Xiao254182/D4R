#!/bin/bash

mkdir -p /opt/d4r
wget https://github.com/Xiao254182/D4R/raw/refs/heads/master/d4r.tar.gz
tar -zxf d4r.tar.gz -C /opt/d4r
chmod +x /opt/d4r/d4r
cp /opt/d4r/d4r /usr/local/bin
echo "使用d4r进入系统"
