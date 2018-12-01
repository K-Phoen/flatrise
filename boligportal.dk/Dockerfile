FROM golang:1.11-stretch

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/go/bin/boligportal.dk"]
