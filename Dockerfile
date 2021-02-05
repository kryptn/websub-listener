FROM golang:1.15.8 AS builder

WORKDIR $GOPATH/src/kryptn/websub-listener/
COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/websub .

FROM scratch


COPY --from=builder /go/bin/websub /go/bin/websub
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


ENTRYPOINT ["/go/bin/websub"]