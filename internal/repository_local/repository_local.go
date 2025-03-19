package repositorylocal

import (
	"fmt"
	test "order-server/pkg/api"
	"sync"

	"github.com/google/uuid"
)

type Repositorylocal struct {
	mutex   sync.Mutex
	storage map[string]*test.Order
}

func New() *Repositorylocal {
	return &Repositorylocal{
		storage: make(map[string]*test.Order),
	}
}

func (r *Repositorylocal) Create(item string, quantity int32) string {
	id := uuid.New() // Генерация нового UUID (v4)
	// копирование sync.Mutex — запрещено, так как это приводит к неопределенному поведению и гонке данных.
	// нужно хранить указатель на test.Order, а не сам объект.
	order := &test.Order{
		Id:       id.String(),
		Item:     item,
		Quantity: quantity,
	}

	// Блокируем перед записью
	r.mutex.Lock()
	r.storage[id.String()] = order
	r.mutex.Unlock()

	return id.String()
}

func (r *Repositorylocal) Update(id string, item string, quantity int32) (*test.Order, error) {
	r.mutex.Lock()
	order, isExist := r.storage[id]
	r.mutex.Unlock()

	if !isExist {
		r.mutex.Unlock()
		return nil, fmt.Errorf("order with specified id does not exist")
	}

	r.mutex.Lock()
	order.Item = item
	order.Quantity = quantity
	r.storage[order.Id] = order
	r.mutex.Unlock()

	return order, nil
}

func (r *Repositorylocal) Get(id string) (*test.Order, error) {
	r.mutex.Lock()
	order, isExist := r.storage[id]
	r.mutex.Unlock()

	if !isExist {
		return nil, fmt.Errorf("order with specified id does not exist")
	}

	return order, nil
}

func (r *Repositorylocal) Delete(id string) (bool, error) {
	r.mutex.Lock()
	_, isExist := r.storage[id]
	r.mutex.Unlock()

	if !isExist {
		return false, fmt.Errorf("order with specified id does not exist")
	}

	delete(r.storage, id)

	return true, nil
}

func (r *Repositorylocal) List() []*test.Order {
	orders := make([]*test.Order, 0, len(r.storage))

	r.mutex.Lock()
	for _, value := range r.storage {
		orders = append(orders, value)
	}
	r.mutex.Unlock()

	return orders
}
