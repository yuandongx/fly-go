FROM golang:1.26.0-alpine3.23 AS builder

# 启用 Go Module支持
ENV GO111MODULE=on

# 配置代理
ENV GOPROXY=https://goproxy.cn,direct

COPY . /app

RUN cd /app;\
    pwd;\
    go build -o fly-go /app/cmd

# 基础镜像：最新版Alpine
FROM alpine:3.23
# 安装必要的运行时依赖（根据程序需求调整）
COPY --from=builder /app/fly-go /usr/local/bin/fly-go

COPY build/entry-point.sh /entry-point.sh

COPY config.yaml /app/config.yaml

RUN chmod +x /entry-point.sh

# 暴露端口（根据程序调整：80=nginx，22=sshd）
EXPOSE 8000

# 设置容器启动命令
ENTRYPOINT ["/entry-point.sh"]