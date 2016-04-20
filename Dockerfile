FROM golang:alpine
MAINTAINER Quadrifoglio <clement@dreamvids.fr>

RUN apk add --update git ffmpeg

RUN mkdir -p /app/src/github.com/dreamvids/convert
COPY . /app/src/github.com/dreamvids/convert
ENV GOPATH /app

WORKDIR /app/src/github.com/dreamvids/convert
RUN go get -d -v
RUN go install -v .

CMD ["/app/bin/convert"]
