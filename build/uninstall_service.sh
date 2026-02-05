#!/bin/bash

# Usage: ./uninstall_service.sh [service_name]
SERVICE_NAME=${1:-fly-go}
UNIT_PATH="/etc/systemd/system/${SERVICE_NAME}.service"

if [ "$EUID" -ne 0 ]; then
  echo "请以 root 或 sudo 身份运行此脚本"
  exit 1
fi

echo "停止并禁用服务（如果存在）"
systemctl stop ${SERVICE_NAME}.service 2>/dev/null || true
systemctl disable ${SERVICE_NAME}.service 2>/dev/null || true

if [ -f "$UNIT_PATH" ]; then
  rm -f "$UNIT_PATH"
  systemctl daemon-reload
  systemctl reset-failed
  echo "已移除系统单元 $UNIT_PATH"
else
  echo "系统单元 $UNIT_PATH 不存在，跳过删除"
fi

if [ -f "/usr/local/bin/${SERVICE_NAME}" ]; then
  rm -f "/usr/local/bin/${SERVICE_NAME}"
  echo "已移除 /usr/local/bin/${SERVICE_NAME}"
fi

echo "完成。"
