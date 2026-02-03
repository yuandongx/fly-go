FROM golang:alpine3.23 AS builder

COPY . /app

RUN cd /app;\
    pwd;\
    go build -o fly-go app/cmd/main.go

FROM alpine:3.23

COPY --from=builder /app/fly-go /app/fly-go

RUN chmod +x /app/fly-go

ENTRYPOINT ["/app/fly-go"]