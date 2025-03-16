FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y \
    make \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/* && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \ 
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest 

ENV PATH="${PATH}:$(go env GOPATH)/bin"

WORKDIR /app

COPY . .

RUN make gen-grpc && \
    make gen-grpc-proxy && \
    make tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/order-server ./cmd/order-server/main.go

FROM ubuntu:22.04

WORKDIR /app 

COPY --from=builder /app/bin/order-server /app/bin/order-server
COPY --from=builder /app/config /app/config
COPY --from=builder /app/db /app/db

EXPOSE 8080
EXPOSE 50051

CMD ["/app/bin/order-server"]
# CMD ["bash"]