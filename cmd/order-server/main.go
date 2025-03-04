package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"order-server/internal/service"
	test "order-server/pkg/api"
	"order-server/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go runRest()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
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

func rungRPC() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
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