FROM golang:1.21.6-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o receipt-processor-app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/receipt-processor-app .

EXPOSE 8080

CMD ["./receipt-processor-app"]