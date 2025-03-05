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
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// go runRest()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed load to config", zap.Error(err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.PortGRPC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := service.New(ctx)
	server := grpc.NewServer()

	test.RegisterOrderServiceServer(server, srv)

	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting serve")

	if err := server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
	}
}

func runRest() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := test.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	guid := uuid.New().String()
	ctx = context.WithValue(ctx, logger.RequestID, guid)

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"request", zap.String("method", info.FullMethod),
		zap.Time("request_time", time.Now()))

	return handler(ctx, req)

}
