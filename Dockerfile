FROM golang:1.18-alpine3.16 as build
RUN mkdir -p /app
WORKDIR /app
COPY . .

ARG netrc
ARG release

ENV CGO_ENABLED=0 RELEASE=$release NETRC=$netrc

RUN echo $NETRC | base64 -d > ~/.netrc && \
    apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev && \
    go mod download &&\
    go build -o main cmd/main.go

FROM alpine:3.13 as prod
WORKDIR /

COPY --from=build /app/main /main

EXPOSE 8083
ENTRYPOINT /main
