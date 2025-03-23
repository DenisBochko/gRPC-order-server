package service

import (
	"context"
	"errors"
	test "order-server/pkg/api"
	"order-server/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// type OrderServiceServer interface {
// 	GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error)
// 	UpdateOrder(context.Context, *UpdateOrderRequest) (*UpdateOrderResponse, error)
// 	DeleteOrder(context.Context, *DeleteOrderRequest) (*DeleteOrderResponse, error)
// 	ListOrders(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error)
// 	mustEmbedUnimplementedOrderServiceServer()
// }

type Repository interface {
	Create(item string, quantity int32) (string, error)
	Update(id string, item string, quantity int32) (*test.Order, error)
	Get(id string) (*test.Order, error)
	Delete(id string) (bool, error)
	List() []*test.Order
}

type Service struct {
	test.OrderServiceServer
	ctx context.Context
	Repository
}

func New(ctx context.Context, repo Repository) *Service {
	return &Service{
		ctx:        ctx,
		Repository: repo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, OrderRequest *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	// Создание заказа
	id, err := s.Repository.Create(
		OrderRequest.GetItem(),
		OrderRequest.GetQuantity(),
	)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to create order", zap.Error(err))
		return nil, errors.New("failed to create order")
	}

	return &test.CreateOrderResponse{Id: id}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, OrderRequest *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	// Апдейт заказа
	order, err := s.Repository.Update(
		OrderRequest.GetId(),
		OrderRequest.GetItem(),
		OrderRequest.GetQuantity(),
	)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to update order", zap.Error(err))
		return nil, err
	}

	return &test.UpdateOrderResponse{Order: order}, nil
}

func (s *Service) GetOrder(ctx context.Context, OrderRequest *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	// Получение заказа
	order, err := s.Repository.Get(
		OrderRequest.GetId(),
	)
	if err != nil {
		return nil, err // Ошибка уже обёрнута
	}

	return &test.GetOrderResponse{Order: order}, nil
}

func (s *Service) DeleteOrder(ctx context.Context, OrderRequest *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	// Удаление заказа
	isSuccess, err := s.Repository.Delete(
		OrderRequest.GetId(),
	)
	if err != nil {
		return nil, err // Ошибка уже обёрнута
	}

	return &test.DeleteOrderResponse{Success: isSuccess}, nil
}

func (s *Service) ListOrders(ctx context.Context, OrdersRequest *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	// Получение всех заказов
	orders := s.Repository.List()

	return &test.ListOrdersResponse{Orders: orders}, nil
}

func (s *Service) LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	guid := uuid.New().String()
	ctx = context.WithValue(s.ctx, logger.RequestID, guid)

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"request", zap.String("method", info.FullMethod),
		zap.Time("request_time", time.Now()))

	return handler(ctx, req)
}
