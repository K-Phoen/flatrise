# build stage
FROM golang:1.11-alpine3.8 as build-env

COPY ./blocket.se /build

RUN apk add --no-cache git gcc musl-dev

RUN cd /build && go build -v -o worker

# final stage
FROM alpine:3.8

WORKDIR /app

RUN apk add --no-cache bash ca-certificates
COPY ./docker/ /docker

COPY --from=build-env /build/worker /app/

CMD ["/app/worker"]
