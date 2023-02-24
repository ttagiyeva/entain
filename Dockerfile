FROM golang:1.17.8-alpine3.15 as build

WORKDIR /go/src/entain
COPY . .
ARG netrc
ARG release
ENV CGO_ENABLED=0 RELEASE=$release NETRC=$netrc GOSUMDB=off


RUN echo $NETRC | base64 -d > ~/.netrc &&  \
    apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go clean --modcache && \
    go mod download && \
    go build -o app cmd/main.go

FROM alpine

COPY --from=build /go/src/entain/migrations /migrations
COPY --from=build /go/src/entain/app /usr/local/bin/app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/app"]
