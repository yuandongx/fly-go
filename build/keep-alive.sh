#!/bin/sh

set -e

while true; do
    echo "keep-alive..."
    # 检查fly-go服务是否正常运行
    if [ ! -f /app/fly-go ]; then
        echo "fly-go not found, restarting..."
        exit 1
    fi
    run=$(ps aux | grep fly-go)
    # 检查fly-go进程是否在运行，如果没有则重启服务
    if [ -z "$run" ]; then
        echo "fly-go not running, starting..."
        exec /app/fly-go
        echo $(stat -c %Y /app/fly-go) > /app/fly-go.last_update
    fi
    # 检查fly-go文件是否更新， 如果更新了则重启服务
    if [ -f /app/fly-go ]; then
        echo "fly-go found, checking update..."
        last_update=$(cat /app/fly-go.last_update)
        latest_update=$(stat -c %Y /app/fly-go)
        if [ "$latest_update" -gt "$last_update" ]; then
            echo "fly-go updated, restarting..."
            killall fly-go || true
            exec /app/fly-go
            echo $(stat -c %Y /app/fly-go) > /app/fly-go.last_update
        fi
    fi
    # Sleep for 5 seconds before printing the next message
    sleep 5
done