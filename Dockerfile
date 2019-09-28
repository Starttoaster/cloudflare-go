#Build image
FROM golang:alpine AS builder
ADD . /go/src/cloudflare
RUN cd /go/src/cloudflare && GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cloudflare \
    && mkdir /config

#End image
FROM alpine
LABEL maintainer="Brandon Butler bmbawb@gmail.com"

RUN apk add --no-cache tzdata
ENV TZ=America/Chicago
COPY --from=builder /go/src/cloudflare/cloudflare /go/src/cloudflare/cloudflare
COPY --from=builder /go/src/cloudflare/templates/index.html /go/src/cloudflare/templates/index.html
COPY --from=builder /config /config
VOLUME /config
ENTRYPOINT ["/go/src/cloudflare/cloudflare"]