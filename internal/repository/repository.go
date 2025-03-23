package repository

import (
	"context"
	"fmt"
	test "order-server/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(item string, quantity int32) (string, error) {
	var id string

	err := r.DB.QueryRow(context.Background(), "INSERT INTO orders (item, quantity) VALUES ($1, $2) RETURNING id", item, quantity).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return id, nil
}

func (r *Repository) Update(id string, item string, quantity int32) (*test.Order, error) {
	order := &test.Order{}

	sts, err := r.DB.Exec(context.Background(), "UPDATE orders SET item = $1, quantity = $2 WHERE id = $3", item, quantity, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	if sts.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	// Устанавливаем поля для возврата пользователю
	order.Id = id
	order.Item = item
	order.Quantity = quantity

	return order, nil
}

func (r *Repository) Get(id string) (*test.Order, error) {
	order := &test.Order{}
	var item string
	var quantity int32

	err := r.DB.QueryRow(context.Background(), "SELECT item, quantity FROM orders WHERE id = $1", id).Scan(&item, &quantity)

	if err == pgx.ErrNoRows {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.Id = id
	order.Item = item
	order.Quantity = quantity

	return order, nil
}

func (r *Repository) Delete(id string) (bool, error) {
	sts, err := r.DB.Exec(context.Background(), "DELETE FROM orders WHERE id = $1", id)

	if sts.RowsAffected() == 0 {
		return false, status.Error(codes.NotFound, "order not found")
	}
	if err != nil {
		return false, fmt.Errorf("failed to delete order: %w", err)
	}

	return true, nil
}

func (r *Repository) List() []*test.Order {
	orders := make([]*test.Order, 0, 100)

	rows, err := r.DB.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		return nil
	}

	for rows.Next() {
		order := &test.Order{}
		err := rows.Scan(&order.Id, &order.Item, &order.Quantity)
		if err != nil {
			return nil
		}

		orders = append(orders, order)
	}

	return orders
}
