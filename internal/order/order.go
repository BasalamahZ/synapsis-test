package order

import (
	"context"
	"time"
)

type Service interface {
	// CreateOrder creates a new order and return the
	// created order ID.
	CreateOrder(ctx context.Context, order Order) (int64, error)

	// GetOrderByID returns a order with the given order ID.
	GetOrderByID(ctx context.Context, id int64) (Order, error)
}

type Order struct {
	ID               int64
	UserID           int64
	UserName         string // derived
	UserEmail        string // derived
	ProductID        int64
	ProductName      string // derived
	ProductPrice     int64  // derived
	Quantity         int64
	TotalAmount      int64
	Status           Status
	ResponseMidtrans string
	CreateTime       time.Time
	UpdateTime       time.Time
}

// Status denotes status of a order.
type Status int

// Followings are the known status.
const (
	StatusUnknown    Status = 0
	StatusSettlement Status = 1
	StatusPending    Status = 2
	StatusCancelled  Status = 3
)

var (
	// StatusList is a list of valid status.
	StatusList = map[Status]struct{}{
		StatusSettlement: {},
		StatusPending:    {},
		StatusCancelled:  {},
	}

	// StatusName maps status to it's string representation.
	statusName = map[Status]string{
		StatusSettlement: "settlement",
		StatusPending:    "pending",
		StatusCancelled:  "cancelled",
	}
)

// Value returns int value of a status type.
func (s Status) Value() int {
	return int(s)
}

// String returns string representaion of a status type.
func (s Status) String() string {
	return statusName[s]
}
