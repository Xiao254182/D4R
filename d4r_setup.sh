#!/bin/bash

# 设置变量
GO_VERSION="1.23.1"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
D4R_URL="https://github.com/Xiao254182/D4R/raw/refs/heads/master/d4r.tar.gz"
D4R_DIR="/opt/d4r"
GO_INSTALL_DIR="/usr/local/go"

# 检查 Go 语言包是否已下载
if [ -d "${GO_INSTALL_DIR}" ]; then
    echo "Go 语言包已下载"
else
    # 部署 Go 语言环境
    echo "正在下载 Go 语言..."
    if wget -q https://golang.google.cn/dl/${GO_TAR}; then
        echo "解压 Go 语言..."
        sudo tar -zxf ${GO_TAR} -C /usr/local
        echo "配置 Go 环境变量..."
        export PATH=\$PATH:${GO_INSTALL_DIR}/bin
        export GOPROXY=https://goproxy.io,direct
        export GOPATH=${D4R_DIR}
    else
        echo "下载 Go 语言失败!"
        exit 1
    fi
fi

# 检查 d4r 是否已下载
if [ -d "${D4R_DIR}" ]; then
    echo "d4r 已下载"
else
# 部署 d4r
    echo "创建目录 ${D4R_DIR}..."
    sudo mkdir -p ${D4R_DIR}

    echo "正在下载 d4r..."
    if wget -q ${D4R_URL} -O d4r.tar.gz; then
        echo "解压 d4r..."
        sudo tar -zxf d4r.tar.gz -C ${D4R_DIR}
        cd ${D4R_DIR} || exit
        echo "编译 d4r..."
        GOOS=linux GOARCH=amd64 go build -o d4r > /dev/null 2>&1
        echo "创建快捷命令..."
        {
            echo "cd ${D4R_DIR} && ./d4r"
        } | sudo tee /usr/local/bin/d4r > /dev/null
        sudo chmod +x /usr/local/bin/d4r
        echo "使用 d4r 进入系统"
    else
        echo "下载 d4r 失败!"
        exit 1
    fi
fi
