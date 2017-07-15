FROM golang:alpine as builder

ADD ./main.go /go/src/github.com/cirocosta/devents/main.go
ADD ./lib /go/src/github.com/cirocosta/devents/lib
ADD ./vendor /go/src/github.com/cirocosta/devents/vendor

WORKDIR /go/src/github.com/cirocosta/devents
RUN set -ex && \
  CGO_ENABLED=0 go build -v -a -ldflags '-extldflags "-static"' && \
  mv ./devents /usr/bin/devents

FROM alpine
COPY --from=builder /usr/bin/devents /usr/local/bin/devents
