package repository

import (
	test "order-server/pkg/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(item string, quantity int32) string {
	return ""
}

func (r *Repository) Update(id string, item string, quantity int32) (*test.Order, error) {
	var order *test.Order
	return order, nil
}

func (r *Repository) Get(id string) (*test.Order, error) {
	var order *test.Order
	return order, nil
}

func (r *Repository) Delete(id string) (bool, error) {
	return true, nil
}

func (r *Repository) List() []*test.Order {
	orders := make([]*test.Order, 0, 0)
	return orders
}
