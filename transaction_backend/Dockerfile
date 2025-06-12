
FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend ./cmd/server

CMD ["./transaction_backend"]
