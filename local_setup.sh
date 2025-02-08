#!/bin/bash

# 设置变量
GO_VERSION="1.23.1"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
D4R_URL="https://github.com/Xiao254182/D4R/archive/refs/heads/master.zip"
D4R_DIR="/opt/d4r"
GO_INSTALL_DIR="/usr/local/go"

# 检查是否有golang环境
if command -v go > /dev/null; then
    echo "golang环境已安装"
else
    echo "golang环境未安装,正在安装..."
        # 检查 wget 是否存在
    if ! command -v wget > /dev/null; then
        echo "wget 未安装，请安装 wget 后重试."
        exit 1
    fi
    wget -q https://golang.google.cn/dl/${GO_TAR}
    echo "解压 Go 语言包..."
    sudo tar -zxf ${GO_TAR} -C /usr/local
    echo "配置 Go 环境变量..."
        {
            echo "export PATH=\$PATH:${GO_INSTALL_DIR}/bin"
            echo "export GOPROXY=https://goproxy.io,direct"
            echo "export GOPATH=${D4R_DIR}"
        } | sudo tee /etc/profile.d/go.sh > /dev/null
    source /etc/profile.d/go.sh
    rm -rf ${GO_TAR}
    echo "golang环境已安装"
fi

# 检查 d4r 是否已下载
mkdir -p ${D4R_DIR}
sudo find /* -type d -name "D4R-*" -exec sh -c 'cp -rf "$1/"* /opt/d4r/' _ {} \;
if [ $? -ne 0 ]; then
    echo "未找到D4R安装包！"
    exit 1
fi
echo "d4r 已下载"
find /* -type d -name "D4R-*" -exec rm -rf {} \; > /dev/null 2>&1

cd ${D4R_DIR} && echo "编译 d4r..."
GOOS=linux GOARCH=amd64 go build -o d4r > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "编译 d4r 失败！"
    exit 1
fi
echo "创建快捷命令..."
{
    echo "cd ${D4R_DIR} && ./d4r"
} | sudo tee /usr/local/bin/d4r > /dev/null
sudo chmod +x /usr/local/bin/d4r
echo "使用 d4r 进入系统"
