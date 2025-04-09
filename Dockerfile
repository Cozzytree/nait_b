FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./cmd/app ./cmd/main.go

CMD ["./cmd/app"]
