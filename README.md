# gRPC order-server

### Запуск через make
```
make build 
make start
```
###### Не забыть поднять базу и изменить config/config.yaml

### Запуск через Docker
###### Перед запуском необходимо сделать make build (сгенерировать файлы из order.proto)

```
docker-compose up --build
```
