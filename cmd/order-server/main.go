package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"order-server/internal/config"
	"order-server/internal/service"
	test "order-server/pkg/api"
	"order-server/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	// Конфигурации
	cfg, err := config.NewENV()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed load to config", zap.Error(err))
	}

	grpcAddr := fmt.Sprintf(":%s", cfg.PortGRPC)
	httpAddr := ":8080"

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := service.New(ctx)
	server := grpc.NewServer(grpc.UnaryInterceptor(srv.LoggerInterceptor))

	test.RegisterOrderServiceServer(server, srv)

	// Создаём контекст для gRPC-Gateway
	gatewayCtx, cancel := context.WithCancel(context.Background())

	// Запускаем gRPC-Gateway
	httpServer, err := runGRPCGateway(gatewayCtx, grpcAddr, httpAddr)
	if err != nil {
		log.Fatalf("failed to start gRPC-Gateway: %v", err)
	}

	// Канал для сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC server is running", zap.String("Port", grpcAddr))
		if err := server.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	// log.Println("Shutting down servers gracefully...")

	// Завершаем gRPC-сервер
	server.GracefulStop()
	// log.Println("gRPC server stopped")

	// Завершаем HTTP сервер с тайм-аутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
	fmt.Println("Server Stopped")

	// Завершаем контекст gRPC-Gateway
	cancel()
}

func runGRPCGateway(ctx context.Context, gGRPAddr, httpAddr string) (*http.Server, error) {
	mux := http.NewServeMux()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := test.RegisterOrderServiceHandlerFromEndpoint(ctx, gwMux, gGRPAddr, opts)
	if err != nil {
		return nil, err
	}

	// Регистрируем gRPC-Gateway на основном HTTP mux
	mux.Handle("/", gwMux)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		log.Println("gRPC-Gateway is running on", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve gRPC-Gateway: %v", err)
		}
	}()

	return server, nil
}
