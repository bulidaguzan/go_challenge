FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git postgresql-client

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]