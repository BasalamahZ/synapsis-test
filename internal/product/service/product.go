package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/synapsis-test/internal/product"
)

// GetProductByID returns a product with the given product ID.
func (s *service) GetProductByID(ctx context.Context, id int64) (product.Product, error) {
	// validate id
	if id <= 0 {
		return product.Product{}, product.ErrInvalidProductID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return product.Product{}, err
	}

	// get a product from postgre
	result, err := pgStoreClient.GetProductByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	return result, nil
}

// GetProducts returns list of products that satisfy the given
// filter.
func (s *service) GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error) {
	// get from redis
	resultCache, err := s.redisClient.Get(ctx, "allProducts").Result()
	if err == nil {
		var products []product.Product
		if err := json.Unmarshal([]byte(resultCache), &products); err != nil {
			return nil, err
		}
		return products, nil
	} else if err != redis.Nil {
		return nil, err
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all products from postgre
	result, err := pgStoreClient.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Cache the data for requests
	resultsJSON, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	if err := s.redisClient.Set(ctx, "allProducts", resultsJSON, 12*time.Hour).Err(); err != nil {
		return nil, err
	}

	return result, nil
}
