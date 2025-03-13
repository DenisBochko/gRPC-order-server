FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y \
    make \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="${PATH}:$(go env GOPATH)/bin"

WORKDIR /app

COPY . .

RUN make install && \
    make gen-grpc && \
    make gen-grpc-proxy && \
    make tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/order-server ./cmd/order-server/main.go

FROM ubuntu:22.04

WORKDIR /app 

COPY --from=builder /app/bin/order-server /app/bin/order-server
COPY --from=builder /app/config /app/config

EXPOSE 8080
EXPOSE 50051

CMD ["/app/bin/order-server"]
# CMD ["bash"]