FROM golang:latest

ENV APP_PATH=/go/src/cloudflare-go
ENV CF_EMAIL="example@email.com"
ENV CF_KEY="my_global_api_key"
ENV CF_ZONE="my_zone_id"

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH
ADD . $APP_PATH
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o cloudflare

CMD go version && go test -cover && go test -v