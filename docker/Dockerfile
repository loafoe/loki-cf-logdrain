FROM golang:1.16.5-alpine3.13 as builder
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /lokiproxy
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o /lokiproxy

FROM grafana/promtail:2.2.1 AS promtail

## Build final image
FROM alpine:3.14.0
LABEL maintainer="andy.lo-a-foe@philips.com"
RUN apk add --no-cache supervisor jq curl && rm -rf /tmp/* /var/cache/apk/*
RUN apk add --no-cache yq --repository http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN mkdir -p /sidecars/bin /sidecars/supervisor/conf.d sidecars/etc /prometheus /thanos

COPY --from=builder     /lokiproxy        /sidecars/bin
COPY --from=promtail    /usr/bin/promtail /sidecars/bin/promtail

COPY docker/supervisord_configs/ /sidecars/supervisor/conf.d
COPY docker/scripts/ /sidecars/bin

EXPOSE 8080

COPY docker/supervisord.conf /etc/
CMD ["supervisord", "--nodaemon", "--configuration", "/etc/supervisord.conf"]
