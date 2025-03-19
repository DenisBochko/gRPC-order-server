package service_test

import (
	"context"
	repositorylocal "order-server/internal/repository_local"
	"order-server/internal/service"
	test "order-server/pkg/api"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тест создания заказа
func TestCreateOrder(t *testing.T) {
	repo := repositorylocal.New()
	service := service.New(context.Background(), repo)

	tests := []struct {
		name string
		*test.CreateOrderRequest
		*test.CreateOrderResponse
		wantErr bool
	}{
		{"valid", &test.CreateOrderRequest{Item: "book1", Quantity: 1}, &test.CreateOrderResponse{Id: uuid.New().String()}, false},
		{"invalid", &test.CreateOrderRequest{Item: "book2", Quantity: 2}, &test.CreateOrderResponse{Id: "123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := service.CreateOrder(context.Background(), tt.CreateOrderRequest)

			if tt.wantErr {
				// require - полностью останавливает тест при ошибке
				require.NotEqual(t, len([]byte(resp.Id)), len([]byte(tt.CreateOrderResponse.Id)))
			} else {
				require.Equal(t, len([]byte(resp.Id)), len([]byte(tt.CreateOrderResponse.Id)))
			}
		})
	}
}

// Тест состояния гонки данных
func TestRaceCondition(t *testing.T) {
	repo := repositorylocal.New()
	service := service.New(context.Background(), repo)

	var wg sync.WaitGroup
	numWorkers := 50

	// Создаём заказ
	req := &test.CreateOrderRequest{Item: "book", Quantity: 10}
	resp, err := service.CreateOrder(context.Background(), req)
	// assert - продолжает выполнение теста при ошибке
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	// Пытаемся конкурентно обновить заказ
	for i := 0; i <= numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updateReq := &test.UpdateOrderRequest{Id: resp.Id, Item: "updated book", Quantity: 20}
			_, _ = service.UpdateOrder(context.Background(), updateReq)
		}()
	}
	wg.Wait()

	// Получаме заказ
	getReq := &test.GetOrderRequest{Id: resp.Id}
	getResp, err := service.GetOrder(context.Background(), getReq)
	assert.NoError(t, err)

	assert.Equal(t, getResp.Order.Id, resp.Id)
	assert.Equal(t, getResp.Order.Item, "updated book")
	assert.Equal(t, getResp.Order.Quantity, int32(20))
}
