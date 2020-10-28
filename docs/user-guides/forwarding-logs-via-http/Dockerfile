FROM golang:1.14-stretch

COPY main.go /go/src/main.go
RUN go build /go/src/main.go

CMD ["/main"]