#
# Stage
#

FROM golang:1.16-buster AS stage

ENV GOPATH="/go" \
    CGO_ENABLED=0

COPY . /tmp/build

WORKDIR /tmp/build
RUN set -x \
    && go build -ldflags '-w -extldflags "-static"' -o regeet-api

#
# Image
#

FROM alpine:3.13

COPY --from=stage /tmp/build/regeet-api /usr/bin/regeet-api
COPY configs/ /etc/regeet

LABEL maintainer="Meraj Sahebdar" \
    io.regeet.name="API" \
    io.regeet.vendor="Regeet" \
    io.regeet.vcs-url="https://github.com/regeet/api" \
    version="0.0.1" \
    license="Apache-2.0"

RUN set -x \
    && apk add --no-cache tzdata ca-certificates \
    && chmod +x /usr/bin/regeet-api

EXPOSE 8080
EXPOSE 8022

CMD = ["regeet-api", "run"]
