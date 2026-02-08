#!/bin/bash

# Usage: ./install_service.sh [path/to/binary] [service_name]
BIN_PATH=${1:-/usr/local/bin/fly-go}
SERVICE_NAME=${2:-fly-go}
TEMPLATE="$(dirname "$0")/fly-go.service.template"
UNIT_PATH="/etc/systemd/system/${SERVICE_NAME}.service"

if [ "$EUID" -ne 0 ]; then
  echo "请以 root 或 sudo 身份运行此脚本"
  exit 1
fi

if [ ! -f "$BIN_PATH" ]; then
  echo "找不到二进制文件: $BIN_PATH"
  exit 1
fi

echo "拷贝二进制到 /usr/local/bin/"
cp "$BIN_PATH" /usr/local/bin/${SERVICE_NAME}
chmod +x /usr/local/bin/${SERVICE_NAME}

echo "生成 systemd 单元文件到 $UNIT_PATH"
sed "s|@@EXEC_PATH@@|/usr/local/bin/${SERVICE_NAME}|g" "$TEMPLATE" > "$UNIT_PATH"

echo "重新加载 systemd，启用并启动服务"
systemctl daemon-reload
systemctl enable ${SERVICE_NAME}.service
systemctl start ${SERVICE_NAME}.service

echo "完成。可以使用 'systemctl status ${SERVICE_NAME}.service' 检查状态。"
