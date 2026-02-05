FROM golang:alpine3.23 AS builder

COPY . /app

RUN cd /app;\
    pwd;\
    go build -o fly-go /app/cmd/main.go

FROM alpine:3.23

# Copy built binary to /usr/local/bin
COPY --from=builder /app/fly-go /usr/local/bin/fly-go

# Copy systemd unit template and render ExecStart to /etc/systemd/system
COPY build/fly-go.service.template /etc/systemd/system/fly-go.service.template

RUN chmod +x /usr/local/bin/fly-go \
 && sed "s|@@EXEC_PATH@@|/usr/local/bin/fly-go|g" /etc/systemd/system/fly-go.service.template > /etc/systemd/system/fly-go.service \
 && mkdir -p /etc/systemd/system/multi-user.target.wants \
 && ln -sf /etc/systemd/system/fly-go.service /etc/systemd/system/multi-user.target.wants/fly-go.service

# Note: To run systemd inside this container you must use a systemd-enabled base image
# and start the container with appropriate privileges (eg. --privileged) and PID 1 as /sbin/init.
# This Dockerfile only installs the unit and binary so the service can be enabled if systemd runs.

ENTRYPOINT ["/usr/local/bin/fly-go"]