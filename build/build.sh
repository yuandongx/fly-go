#!/bin/bash

# 检查脚本执行路径是否为项目根目录
pth=$(basename $1)

if [ $pth != "build" ]; then
    echo "请在项目根目录执行脚本"
    cd ..
fi
pth=$(basename $PWD)
if [ $pth != "fly-go" ]; then
    echo "请在项目根目录执行脚本"
    exit 1
fi

echo $PWD
