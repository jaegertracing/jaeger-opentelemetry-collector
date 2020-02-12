FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.13-buster as builder
WORKDIR /go/src/app
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod go.sum Makefile ./
RUN ls -al
RUN go version
# TODO LD_FLAGS are wrong, e.g. there is no git in the image
RUN make build

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/app/cmd/collector/jaeger-opentelemetry-collector /jaeger-opentelemetry-collector
ENTRYPOINT ["/jaeger-opentelemetry-collector"]
EXPOSE 55678 55679
