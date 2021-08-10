FROM golang:1.13.6-alpine3.11 as buildergo
RUN mkdir -p /fluent-bit
COPY cmd/fluent-bit-watcher/main.go go.mod /fluent-bit/
WORKDIR /fluent-bit
RUN CGO_ENABLED=0 go build -i -ldflags '-w -s' -o fluent-bit main.go

# FROM gcr.io/distroless/cc-debian10
FROM fluent/fluent-bit:1.8.3
MAINTAINER KubeSphere <kubesphere@yunify.com>
LABEL Description="Fluent Bit docker image" Vendor="KubeSphere" Version="1.0"

COPY conf/fluent-bit.conf conf/parsers.conf /fluent-bit/etc/
COPY --from=buildergo /fluent-bit/fluent-bit /fluent-bit/bin/fluent-bit-watcher

#
EXPOSE 2020

# Entry point
CMD ["/fluent-bit/bin/fluent-bit-watcher", "-c", "/fluent-bit/etc/fluent-bit.conf"]
