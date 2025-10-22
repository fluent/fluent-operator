ARG FLUENT_BIT_BASE_VERSION=3.2.4 \
    GO_VERSION=1.24.5

FROM golang:${GO_VERSION}-alpine3.21 AS buildergo
RUN mkdir -p /fluent-bit
RUN mkdir -p /code
COPY . /code/
WORKDIR /code
RUN echo $(ls -al /code)
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /fluent-bit/fluent-bit /code/cmd/fluent-watcher/fluentbit/main.go

FROM fluent/fluent-bit:${FLUENT_BIT_BASE_VERSION}
LABEL Description="Fluent Bit docker image" Vendor="Fluent" Version="1.0"

COPY conf/fluent-bit.conf conf/fluent-bit.yaml conf/parsers.conf /fluent-bit/etc/
COPY --from=buildergo /fluent-bit/fluent-bit /fluent-bit/bin/fluent-bit-watcher

# Entry point
ENTRYPOINT ["/fluent-bit/bin/fluent-bit-watcher"]
