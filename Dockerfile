FROM golang:1.21.6-alpine3.18 AS builder
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite build-base
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN CGO_ENABLED=1 go build \
    -o /build/http \
    -ldflags="-w -s" \
    -gcflags="all=-c=2" \
    ./main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/http /app/http
EXPOSE 8080
CMD ["/app/http"]
