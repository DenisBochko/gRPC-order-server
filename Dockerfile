FROM golang:1.24 AS builder

WORKDIR /app 

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/order-server ./cmd/order-server/main.go

FROM ubuntu:22.04

WORKDIR /app 

COPY --from=builder /app/bin/order-server /app/bin/order-server
COPY --from=builder /app/config /app/config

EXPOSE 8080
EXPOSE 50051

CMD ["/app/bin/order-server"]