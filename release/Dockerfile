#Build image
FROM golang:latest AS builder

ENV APP_PATH=/go/src/cloudflare-go

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

ADD . $APP_PATH
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o cloudflare


#End image
FROM alpine
LABEL maintainer="Brandon Butler bmbawb@gmail.com"

RUN apk add --no-cache tzdata
ENV TZ=America/Chicago
COPY --from=builder /go/src/cloudflare-go/cloudflare /cloudflare
COPY --from=builder /go/src/cloudflare-go/index.html /index.html

ENTRYPOINT ["/cloudflare"]