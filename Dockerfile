FROM golang:1.12.1-alpine as builder

RUN apk update && apk add git && apk add ca-certificates && adduser -u 1000 -D app

WORKDIR /webhook-server

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/webhook-server cmd/server/main.go


FROM scratch AS base
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/webhook-server /bin/webhook-server
USER app
ENTRYPOINT ["/bin/webhook-server", "-L", "debug"]
