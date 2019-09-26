#Build image
FROM golang:alpine AS builder
COPY ./cloudflare.go /go/cloudflare.go
RUN cd /go && GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cloudflare \
    && mkdir /config

#End image
FROM alpine
LABEL maintainer="Brandon Butler bmbawb@gmail.com"

COPY --from=builder /go/cloudflare /go/cloudflare
COPY --from=builder /config /config
VOLUME /config
ENTRYPOINT ["/go/cloudflare"]