FROM golang:1.11-stretch

WORKDIR /app

COPY . .

RUN go build -v -o worker

CMD ["/app/worker"]
