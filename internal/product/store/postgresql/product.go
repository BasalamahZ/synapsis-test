package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/product"
)

// GetProductByID returns a product with the given product ID.
func (sc *storeClient) GetProductByID(ctx context.Context, id int64) (product.Product, error) {
	query := fmt.Sprintf(queryGetProduct, "WHERE p.id = $1")
	// query single row
	var pdb productDB
	err := sc.q.QueryRowx(query, id).StructScan(&pdb)
	if err != nil {
		if err == sql.ErrNoRows {
			return product.Product{}, product.ErrDataNotFound
		}
		return product.Product{}, err
	}

	return pdb.format(), nil
}

// GetProducts returns list of products that satisfy the given
// filter.
func (sc *storeClient) GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if filter.CategoryID > 0 {
		addConditions = append(addConditions, "p.category_id = :category_id")
		argsKV["category_id"] = filter.CategoryID
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}

	// construct query
	query := fmt.Sprintf(queryGetProduct, addCondition)

	// prepare query
	query, args, err := sqlx.Named(query, argsKV)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sc.q.Rebind(query)

	// query to database
	rows, err := sc.q.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// read products
	products := make([]product.Product, 0)
	for rows.Next() {
		var row productDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		products = append(products, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
