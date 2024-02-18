package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/synapsis-test/internal/order"
)

func (s *service) CreateOrder(ctx context.Context, reqOrder order.Order) (int64, error) {
	// validate fields
	err := validateOrder(reqOrder)
	if err != nil {
		return 0, err
	}

	carts, err := s.product.GetCartsByUserID(ctx, reqOrder.UserID)
	if err != nil {
		return 0, err
	}

	// validate order product in cart
	valid := false
	for _, cart := range carts {
		fmt.Println("cart.ProductID", cart.ProductID)
		fmt.Println("reqOrder.ProductID", reqOrder.ProductID)
		if cart.ProductID == reqOrder.ProductID {
			valid = true
			reqOrder.TotalAmount = reqOrder.Quantity * cart.ProductPrice
			reqOrder.ProductName = cart.ProductName
			reqOrder.ProductPrice = cart.ProductPrice
			break
		}

	}
	fmt.Println("valid", valid)
	if !valid {
		return 0, order.ErrInvalidOrderID
	}

	// update fields
	reqOrder.CreateTime = s.timeNow()

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	resPayment, err := s.payment(reqOrder)
	if err != nil {
		return 0, err
	}

	jsonData, err := json.Marshal(resPayment)
	if err != nil {
		return 0, err
	}

	reqOrder.ResponseMidtrans = string(jsonData)

	// create order in pgstore
	orderID, err := pgStoreClient.CreateOrder(ctx, reqOrder)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (s *service) GetOrderByID(ctx context.Context, id int64) (order.Order, error) {
	// validate id
	if id <= 0 {
		return order.Order{}, order.ErrInvalidOrderID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return order.Order{}, err
	}

	// get a product from postgre
	result, err := pgStoreClient.GetOrderByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}

	return result, nil
}

// validateOrder validates fields of the given
// order.
func validateOrder(reqOrder order.Order) error {
	if reqOrder.UserID <= 0 {
		return order.ErrInvalidUserID
	}

	if reqOrder.ProductID <= 0 {
		return order.ErrInvalidProductID
	}

	if reqOrder.Quantity <= 0 {
		return order.ErrInvalidQuantity
	}

	return nil
}
