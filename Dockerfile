FROM golang:1-alpine
LABEL maintainer="Brandon Butler bmbawb@gmail.com"

WORKDIR /go
COPY ./cloudflare.go /go/cloudflare.go

RUN cd /go && go build -o /go/cloudflare

VOLUME /go

CMD /go/cloudflare
