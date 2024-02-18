package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/product"
)

func (sc *storeClient) AddProductCart(ctx context.Context, reqProductCart product.ProductCart) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"user_id":     reqProductCart.UserID,
		"product_id":  reqProductCart.ProductID,
		"quantity":    reqProductCart.Quantity,
		"create_time": reqProductCart.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryAddProductCart, argsKV)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = sc.q.Rebind(query)

	// execute query
	_, err = sc.q.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (sc *storeClient) GetCartsByUserID(ctx context.Context, userID int64) ([]product.ProductCart, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if userID > 0 {
		addConditions = append(addConditions, "pc.user_id = :user_id")
		argsKV["user_id"] = userID
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}

	// construct query
	query := fmt.Sprintf(queryGetProductCart, addCondition)

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

	// read product cart
	productCart := make([]product.ProductCart, 0)
	for rows.Next() {
		var row productCartDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		productCart = append(productCart, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productCart, nil
}

func (sc *storeClient) DeleteProductCartByID(ctx context.Context, reqProductCart product.ProductCart) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"user_id":    reqProductCart.UserID,
		"product_id": reqProductCart.ProductID,
	}

	// prepare query
	query, args, err := sqlx.Named(queryDeleteProductCart, argsKV)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = sc.q.Rebind(query)

	// execute query
	_, err = sc.q.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
