install: # Установка плагинов для генерации кода 
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen-grpc: # Генерация grpc из proto файла
	protoc --go_out=./pkg/api/ --go_opt=paths=source_relative --go-grpc_out=./pkg/api/ --go-grpc_opt=paths=source_relative order.proto

gen-grpc-proxy: # Генерация grpc-gateway из proto файла
	protoc --grpc-gateway_out=allow_delete_body=true:./pkg/api/ --grpc-gateway_opt paths=source_relative order.proto

tidy: # Установка зависимостей
	go mod tidy

build: # Сборка бинарника
	go build -o ./bin/ ./cmd/order-server

start: # Старт 
	./bin/order-server
