# gRPC order-server

### Запуск через make
```
make build 
make start
```

### Запуск через Docker

```
docker build -t grpc-server:1 -f Dockerfile.dev .
docker run --rm -p 8080:8080 -p 50051:50051 --name grpc-server grpc-server:1 
```
