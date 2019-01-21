FROM alpine:3.6

RUN adduser -D fluentbit-operator
USER fluentbit-operator

ADD tmp/_output/bin/fluentbit-operator /usr/local/bin/fluentbit-operator
