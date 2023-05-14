# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.20.3-alpine3.17 as builder

ARG VERSION

ADD ./ /source
RUN cd /source && \
addgroup -g 1002 appuser && \
adduser -H -u 1002 -D -G appuser appuser && \
go build -ldflags="-X 'main.Version=${VERSION}'" -o spectromate -v

FROM alpine:latest

LABEL org.opencontainers.image.source="http://spectrocloud.com/spectromate"
LABEL org.opencontainers.image.description "Spectromate is an API server with features to support Slack bot integration."


COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appuser  /source/spectromate /usr/bin/
COPY entrypoint.sh /usr/bin/

RUN apk -U upgrade && apk add bash jq git --no-cache && mkdir /packs && chown -R appuser:appuser /packs && \
mkdir -p /var/log/spectromate && chown -R appuser:appuser /var/log/spectromate && \
touch /var/log/spectromate.log && chown appuser:appuser /var/log/spectromate.log && chmod 664 /var/log/spectromate.log
USER appuser

CMD ["/usr/bin/entrypoint.sh"]
