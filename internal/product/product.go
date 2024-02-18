package product

import (
	"context"
	"time"
)

type Service interface {
	// GetProductByID returns a product with the given product ID.
	GetProductByID(ctx context.Context, id int64) (Product, error)

	// GetProducts returns list of products that satisfy the given
	// filter.
	GetProducts(ctx context.Context, filter GetProductsFilter) ([]Product, error)

	// AddProductCart add a prodcut to cart
	AddProductCart(ctx context.Context, productCart ProductCart) error

	// GetCartsByUserID returns all product in cart with the given user ID.
	GetCartsByUserID(ctx context.Context, userID int64) ([]ProductCart, error)

	// DeleteProductCartByID deletes a product in cart with the given cart ID.
	DeleteProductCartByID(ctx context.Context, productCart ProductCart) error
}

type Product struct {
	ID           int64
	Name         string
	Description  string
	Price        int64
	CategoryID   int64
	CategoryName string //derived
	CreateTime   time.Time
	UpdateTime   time.Time
}

type GetProductsFilter struct {
	CategoryID int64
}

type ProductCart struct {
	UserID       int64
	ProductID    int64
	ProductName  string // derived
	ProductPrice int64  // derived
	Quantity     int64
	CreateTime   time.Time
	UpdateTime   time.Time
}
