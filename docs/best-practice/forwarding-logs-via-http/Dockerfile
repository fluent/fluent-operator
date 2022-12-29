# Build the manager binary
FROM golang:1.19.4 as builder

WORKDIR /
COPY main.go /go/src/main.go
RUN CGO_ENABLED=0 go build /go/src/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM kubesphere/distroless-static:nonroot
WORKDIR /
COPY --from=builder /main .
USER nonroot:nonroot

ENTRYPOINT ["/main"]

