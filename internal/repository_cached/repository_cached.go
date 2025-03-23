package repositorycached

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	test "order-server/pkg/api"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Create(item string, quantity int32) (string, error)
	Update(id string, item string, quantity int32) (*test.Order, error)
	Get(id string) (*test.Order, error)
	Delete(id string) (bool, error)
	List() []*test.Order
}

type RepositoryCached struct {
	repo  Repository
	redis *redis.Client
	ttl   time.Duration
}

func New(repo Repository, redis *redis.Client, ttl time.Duration) *RepositoryCached {
	return &RepositoryCached{
		repo:  repo,
		redis: redis,
		ttl:   ttl,
	}
}

func (r *RepositoryCached) Create(item string, quantity int32) (string, error) {
	order := &test.Order{}

	id, err := r.repo.Create(item, quantity)
	if err != nil {
		return "", err
	}

	// Устанавливаем поля для возврата пользователю
	order.Id = id
	order.Item = item
	order.Quantity = quantity

	// После создания запишем в кэш
	data, err := json.Marshal(order)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order: %w", err)
	}

	err = r.redis.Set(context.Background(), order.Id, data, r.ttl).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set order to cache: %w", err)
	}

	log.Println("Order created and set to cache")

	return id, nil
}

func (r *RepositoryCached) Update(id string, item string, quantity int32) (*test.Order, error) {
	order := &test.Order{}

	// Обновляем заказ в БД
	order, err := r.repo.Update(id, item, quantity)
	if err != nil {
		return nil, err
	}

	// Обновляем кэш
	data, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order: %w", err)
	}

	err = r.redis.Set(context.Background(), order.Id, data, r.ttl).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set order to cache: %w", err)
	}

	log.Println("Order updated and set to cache")

	return order, nil
}

func (r *RepositoryCached) Get(id string) (*test.Order, error) {
	order := &test.Order{}

	// Пытаемся получить заказ из кэша
	data, err := r.redis.Get(context.Background(), id).Result()
	if err == redis.Nil {
		order, err := r.repo.Get(id)
		if err != nil {
			return nil, err
		}

		log.Println("Order from DB")

		// Обновляем кэш
		updatingData, err := json.Marshal(order)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal order: %w", err)
		}

		err = r.redis.Set(context.Background(), order.Id, updatingData, r.ttl).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to set order to cache: %w", err)
		}

		return order, nil
	} else if err != nil {
		return nil, err
	}

	// Если заказ найден в кэше, то возвращаем его
	err = json.Unmarshal([]byte(data), order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	log.Println("Order from cache")

	return order, nil
}

func (r *RepositoryCached) Delete(id string) (bool, error) {
	// Удаляем заказ из БД
	deleted, err := r.repo.Delete(id)
	if err != nil {
		return false, err
	}

	// Удаляем заказ из кэша
	err = r.redis.Del(context.Background(), id).Err()
	if err != nil {
		return false, fmt.Errorf("failed to delete order from cache: %w", err)
	}

	return deleted, nil
}

func (r *RepositoryCached) List() []*test.Order {
	return r.repo.List()
}
