FROM golang:1.20.3-alpine3.17 as builder

LABEL org.opencontainers.image.source="http://spectrocloud.com/spectromate"
LABEL org.opencontainers.image.description "An API server with features to support Slack bots integration."

ARG VERSION

ADD ./ /source
RUN cd /source && \
adduser -H -u 1002 -D appuser appuser && \
go build -ldflags="-X 'spectrocloud.com/spectromate/cmd.VersionString=${VERSION}'" -o spectromate -v

FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appuser  /source/spectromate /usr/bin/

RUN apk -U upgrade && apk add bash jq git --no-cache && mkdir /packs -p && chown appuser:appuser /packs
USER appuser

ENTRYPOINT ["/usr/bin/spectromate"]