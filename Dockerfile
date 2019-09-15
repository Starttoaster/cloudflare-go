FROM golang:1-alpine
LABEL maintainer="Brandon Butler bmbawb@gmail.com"

WORKDIR /go
COPY ./cloudflare.go /go/cloudflare.go

RUN cd /go && go build -o /go/cloudflare \
    && mkdir /config

VOLUME /config

CMD /go/cloudflare
