FROM golang:alpine as builder

ADD ./main.go /go/src/github.com/cirocosta/devents/main.go
ADD ./lig /go/src/github.com/cirocosta/devents/lib
ADD ./vendor /go/src/github.com/cirocosta/devents/vendor

WORKDIR /go/src/github.com/cirocosta/devents/devents
RUN set -ex && \
  CGO_ENABLED=0 go build -v -a -ldflags '-extldflags "-static"' && \
  mv ./devents /usr/bin/devents

FROM busybox
COPY --from=builder /usr/bin/devents /devents

CMD [ "devents" ]
