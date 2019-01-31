FROM golang:alpine
MAINTAINER KAMEIKE

RUN apk add --update gcc musl-dev
RUN apk add --update git

ADD . /go/src/github.com/kameike/karimono
WORKDIR /go/src/github.com/kameike/karimono
RUN go get .

RUN apk add --update sqlite
RUN apk add --update sqlite-dev

RUN go build --tags "libsqlite3 linux"

FROM alpine 
# RUN apk add --update gcc musl-dev
COPY --from=0 /go/bin/karimono .
ENV PORT 8080
CMD ["./karimono"]
