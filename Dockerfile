FROM golang:1.21 as build

RUN mkdir -p /app
WORKDIR /app
COPY .. .
RUN go mod download && \
    go build -o main cmd/main.go

FROM alpine:3.16 as prod

WORKDIR /
RUN apk add libc6-compat
COPY --from=build /app/main /main

EXPOSE 8080

ENTRYPOINT /main
