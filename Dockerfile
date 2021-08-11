#
# Stage
#

FROM golang:1.16-buster AS stage

ENV GOPATH="/go" \
    CGO_ENABLED=0

COPY . /tmp/build

WORKDIR /tmp/build
RUN set -x \
    && go build -ldflags '-w -extldflags "-static"' -o bitban-api

#
# Image
#

FROM alpine:3.13

COPY --from=stage /tmp/build/bitban-api /usr/bin/bitban-api
COPY configs/ /etc/bitban

LABEL maintainer="Meraj Sahebdar" \
    io.bitban.name="server" \
    io.bitban.vendor="bitban" \
    io.bitban.vcs-url="https://github.com/merajsahebdar/bitban-server" \
    version="0.0.1" \
    license="Apache-2.0"

RUN set -x \
    && apk add --no-cache tzdata ca-certificates \
    && chmod +x /usr/bin/bitban-api

EXPOSE 8080
EXPOSE 8022

CMD = ["bitban-api", "run"]
