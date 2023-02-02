FROM golang:1.19.5-buster as builder
WORKDIR /loki-cf-logdrain
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o /loki-cf-logdrain

FROM grafana/promtail:2.7.3 AS promtail

## Build final image
FROM debian:buster-slim
LABEL maintainer="andy.lo-a-foe@philips.com"
RUN apt-get update && apt-get install -y \
    jq \
    supervisor \
    curl \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /sidecars/bin /sidecars/supervisor/conf.d sidecars/etc /promtail

COPY --from=builder     /loki-cf-logdrain        /sidecars/bin
COPY --from=promtail    /usr/bin/promtail /sidecars/bin/promtail

COPY docker/supervisord_configs/ /sidecars/supervisor/conf.d
COPY docker/scripts/ /sidecars/bin

EXPOSE 8080

COPY docker/supervisord.conf /etc/
CMD ["supervisord", "--nodaemon", "--configuration", "/etc/supervisord.conf"]
