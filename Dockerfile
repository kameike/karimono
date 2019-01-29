FROM golang
MAINTAINER KAMEIKE
RUN go get -v -u github.com/kameike/karimono

FROM alpine
COPY --from=0 /go/bin/karimono .
ENV PORT 8080
CMD ["./karimono"]
