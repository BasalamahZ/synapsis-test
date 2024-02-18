package service

import (
	"context"

	"github.com/synapsis-test/internal/product"
)

func (s *service) AddProductCart(ctx context.Context, reqProductCart product.ProductCart) error {
	// validate fields
	err := validateProductCart(reqProductCart)
	if err != nil {
		return err
	}

	// update fields
	reqProductCart.CreateTime = s.timeNow()

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// add product cart in pgstore
	err = pgStoreClient.AddProductCart(ctx, reqProductCart)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetCartsByUserID(ctx context.Context, userId int64) ([]product.ProductCart, error) {
	// validate id
	if userId <= 0 {
		return nil, product.ErrInvalidUserID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all product carts from postgre
	result, err := pgStoreClient.GetCartsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) DeleteProductCartByID(ctx context.Context, reqProductCart product.ProductCart) error {
	if reqProductCart.UserID <= 0 {
		return product.ErrInvalidUserID
	}

	if reqProductCart.ProductID <= 0 {
		return product.ErrInvalidProductID
	}

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// delete product cart in pgstore
	err = pgStoreClient.DeleteProductCartByID(ctx, reqProductCart)
	if err != nil {
		return err
	}

	return nil
}

// validateProductCart validates fields of the given
// product cart.
func validateProductCart(reqProductCart product.ProductCart) error {
	if reqProductCart.UserID <= 0 {
		return product.ErrInvalidUserID
	}

	if reqProductCart.ProductID <= 0 {
		return product.ErrInvalidProductID
	}

	if reqProductCart.Quantity <= 0 {
		return product.ErrInvalidQuantity
	}

	return nil
}
