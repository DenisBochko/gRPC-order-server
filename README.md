# gRPC order-server

### Запуск через make
##### Не забыть поднять базу и изменить config/config.yaml
```
make setup 
make start
```

### Запуск через Docker
#### Необходимо создать ```.env``` файл (есть пример в env.example)
```
docker-compose up --build
```

### Тестирование service
```
go test -count=10 --race ./internal/service
```