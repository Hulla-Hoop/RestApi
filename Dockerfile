FROM golang:latest as builder
WORKDIR /app
COPY . /app
RUN GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o app cmd/app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY ./.env /app
RUN mkdir -p /app/migrations
COPY ./migrations/  /app/migrations
ENTRYPOINT ["./app"]
