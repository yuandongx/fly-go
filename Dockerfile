FROM golang:1.19-alpine AS builder

COPY . /app

RUN cd /app && go build -o fly-go

COPY --from=builder /app/fly-go /app/fly-go

