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
	"order-server/pkg/postgres"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	// Обработка сигналов завершения
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Конфигурации
	cfg, err := config.NewENV()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed load to config", zap.Error(err))
	}

	// Подключение к БД
	conn, err := postgres.New(ctx, cfg.PostgresCfg)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to connect to database", zap.Error(err))
	}

	if conn.Ping(ctx) != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
	}

	// Порты
	grpcAddr := fmt.Sprintf(":%s", cfg.PortGRPC)
	httpAddr := fmt.Sprintf(":%s", cfg.PortHttp)

	// Запускаем gRPC-Gateway
	httpServer, err := runGRPCGateway(ctx, grpcAddr, httpAddr)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to start gRPC-Gateway", zap.Error(err))
	}

	// Запускаем gRPC server
	grpcServer, err := runGRPC(ctx, grpcAddr)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to start gRPC server", zap.Error(err))
	}

	// Ожидаем сигнал завершения
	select {
	case <-ctx.Done():
		log.Println("Shutting down servers gracefully...")

		// Завершаем gRPC-сервер
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")

		// Завершаем Pool соединений с БД
		conn.Close()
		log.Println("Database connection closed")

		// Завершаем HTTP сервер
		httpServer.Shutdown(ctx)
		log.Println("Http server stopped")

		log.Println("Server Stopped")
	}
}

func runGRPCGateway(ctx context.Context, gGRPAddr, httpAddr string) (*http.Server, error) {
	mux := http.NewServeMux()

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
		logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC-Gateway is running", zap.String("Port", httpAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve gRPC-Gateway: %v", err)
		}
	}()

	return server, nil
}

func runGRPC(ctx context.Context, gGRPAddr string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", gGRPAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen TCP: %e", err)
	}

	srv := service.New(ctx)
	server := grpc.NewServer(grpc.UnaryInterceptor(srv.LoggerInterceptor))

	test.RegisterOrderServiceServer(server, srv)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC server is running", zap.String("Port", gGRPAddr))
		if err := server.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
		}
	}()

	return server, nil
}
