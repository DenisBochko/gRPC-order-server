package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"order-server/internal/service"
	test "order-server/pkg/api"
	"order-server/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
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
