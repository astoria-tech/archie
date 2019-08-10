FROM golang:1.12-alpine as builder

RUN mkdir /build
ADD . /build

WORKDIR /build
RUN go build -ldflags="-w -s" -mod vendor -o /go/bin/main

FROM alpine

RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /go/bin/main /app/

VOLUME /app/config.yaml

WORKDIR /app
CMD ["./main"]
