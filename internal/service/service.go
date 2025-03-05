package service

import (
	"context"
	"fmt"
	test "order-server/pkg/api"
	"order-server/pkg/logger"
	"sync"
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

type Service struct {
	test.OrderServiceServer
	ctx     context.Context
	mutex   sync.Mutex
	storage map[string]*test.Order
}

func New(ctx context.Context) *Service {
	return &Service{
		ctx:     ctx,
		storage: make(map[string]*test.Order),
	}
}

/*
	{
	    "item": "book1",
	    "quantity": 122
	}
*/
func (s *Service) CreateOrder(ctx context.Context, OrderRequest *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	id := uuid.New() // Генерация нового UUID (v4)

	// копирование sync.Mutex — запрещено, так как это приводит к неопределенному поведению и гонке данных.
	// нужно хранить указатель на test.Order, а не сам объект.
	order := &test.Order{
		Id:       id.String(),
		Item:     OrderRequest.GetItem(),
		Quantity: OrderRequest.GetQuantity(),
	}

	// Блокируем перед записью
	s.mutex.Lock()
	s.storage[id.String()] = order
	s.mutex.Unlock()

	// Логируем создание заказа
	// logger.GetLoggerFromCtx(s.ctx).Info(s.ctx, fmt.Sprintf("created order: %v", order))

	return &test.CreateOrderResponse{Id: id.String()}, nil
}

/*
	{
		"id": "",
	    "item": "book1",
	    "quantity": 122
	}
*/
func (s *Service) UpdateOrder(ctx context.Context, OrderRequest *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	id := OrderRequest.GetId()

	s.mutex.Lock()
	order, isExist := s.storage[id]
	s.mutex.Unlock()

	if !isExist {
		return nil, fmt.Errorf("order with specified id does not exist")
	}

	order.Item = OrderRequest.GetItem()
	order.Quantity = OrderRequest.GetQuantity()

	s.mutex.Lock()
	s.storage[order.Id] = order
	s.mutex.Unlock()

	// logger.GetLoggerFromCtx(s.ctx).Info(s.ctx, fmt.Sprintf("update order: %v", order))

	return &test.UpdateOrderResponse{Order: order}, nil
}

/*
	{
	   "id": ""
	}
*/
func (s *Service) GetOrder(ctx context.Context, OrderRequest *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	id := OrderRequest.GetId()

	s.mutex.Lock()
	order, isExist := s.storage[id]
	s.mutex.Unlock()

	if !isExist {
		return nil, fmt.Errorf("order with specified id does not exist")
	}

	// logger.GetLoggerFromCtx(s.ctx).Info(s.ctx, fmt.Sprintf("geted order: %v", order))

	return &test.GetOrderResponse{Order: order}, nil
}

/*
	{
	   "id": ""
	}
*/
func (s *Service) DeleteOrder(ctx context.Context, OrderRequest *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	id := OrderRequest.GetId()

	s.mutex.Lock()
	_, isExist := s.storage[id]
	s.mutex.Unlock()

	if !isExist {
		return &test.DeleteOrderResponse{Success: false}, fmt.Errorf("order with specified id does not exist")
	}

	delete(s.storage, id)

	// logger.GetLoggerFromCtx(s.ctx).Info(s.ctx, fmt.Sprintf("deleted order with id: %s", id))

	return &test.DeleteOrderResponse{Success: true}, nil
}

func (s *Service) ListOrders(ctx context.Context, OrdersRequest *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	orders := make([]*test.Order, 0, len(s.storage))

	s.mutex.Lock()
	for _, value := range s.storage {
		orders = append(orders, value)
	}
	s.mutex.Unlock()

	// logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint(ctx.Value(logger.RequestID)))

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
