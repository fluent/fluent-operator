FROM golang:1.10-alpine as golang

ADD . /go/src/kubesphere.io/fluentbit-operator
WORKDIR /go/src/kubesphere.io/fluentbit-operator

RUN apk add --update --no-cache ca-certificates curl git make
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -v -vendor-only

RUN go install ./cmd/fluentbit-operator


FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=golang /go/bin/fluentbit-operator /usr/local/bin/fluentbit-operator

RUN adduser -D fluentbit-operator
USER fluentbit-operator

ENTRYPOINT ["/usr/local/bin/fluentbit-operator"]
