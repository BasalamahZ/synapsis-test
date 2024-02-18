package order

import "errors"

// Followings are the known errors returned from order.
var (
	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidOrderID is returned when the given order ID is
	// invalid.
	ErrInvalidOrderID = errors.New("invalid order id")

	// ErrInvalidProductID is returned when the given product ID is
	// invalid.
	ErrInvalidProductID = errors.New("invalid product id")

	// ErrInvalidUserID is returned when the given user ID is
	// invalid.
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidQuantity is returned when the given quantity is
	// invalid.
	ErrInvalidQuantity = errors.New("invalid quantity")
)
