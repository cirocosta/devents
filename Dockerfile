FROM golang:alpine as builder

ADD ./vendor /go/src/github.com/cirocosta/devents/vendor
ADD ./devents /go/src/github.com/cirocosta/devents/devents

WORKDIR /go/src/github.com/cirocosta/devents/devents
RUN set -ex && \
  CGO_ENABLED=0 go build -v -a -ldflags '-extldflags "-static"' && \
  mv ./devents /usr/bin/devents

FROM busybox
COPY --from=builder /usr/bin/devents /devents

CMD [ "devents" ]
