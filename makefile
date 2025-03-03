install: # Установка плагинов для генерации кода 
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen: # Генерация из proto файла
	protoc --go_out=./pkg/api/ --go_opt=paths=source_relative --go-grpc_out=./pkg/api/ --go-grpc_opt=paths=source_relative order.proto

tidy: # Установка зависимостей
	go mod tidy

build: # Сборка бинарника
	go build -o ./bin/ ./cmd/order-server

start: # Старт 
	./bin/order-server
