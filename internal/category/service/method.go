package service

import (
	"context"

	"github.com/synapsis-test/internal/category"
)

func (s *service) GetCategoryByID(ctx context.Context, id int64) (category.Category, error) {
	// validate id
	if id <= 0 {
		return category.Category{}, category.ErrInvalidCategoryID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return category.Category{}, err
	}

	// get a product from postgre
	result, err := pgStoreClient.GetCategoryByID(ctx, id)
	if err != nil {
		return category.Category{}, err
	}

	return result, nil
}

func (s *service) GetCategories(ctx context.Context) ([]category.Category, error) {
	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all product from postgre
	result, err := pgStoreClient.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
