FROM golang:latest AS builder

COPY . /app

RUN cd /app;\
    pwd;\
    go build -o fly-go /app/cmd/main.go

# 基础镜像：最新版Alpine
FROM alpine:latest

# Copy built binary to /usr/local/bin
COPY --from=builder /app/fly-go /usr/local/bin/fly-go

# 更新源并安装依赖（supervisord + 示例程序：nginx、openssh-server）
# nginx \
# openssh-server \
RUN apk update && apk add --no-cache supervisor \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /etc/supervisor/conf.d /var/log/supervisor

# 复制Supervisord主配置文件到容器
COPY --from=builder /app/build/supervisord.conf /etc/supervisor/supervisord.conf

# 暴露端口（根据程序调整：80=nginx，22=sshd）
EXPOSE 8000

# 启动Supervisord（前台运行，作为容器主进程）
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]