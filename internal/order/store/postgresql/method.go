package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/order"
)

func (sc *storeClient) CreateOrder(ctx context.Context, reqOrder order.Order) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"user_id":           reqOrder.UserID,
		"product_id":        reqOrder.ProductID,
		"quantity":          reqOrder.Quantity,
		"total_amount":      reqOrder.TotalAmount,
		"status":            reqOrder.Status,
		"response_midtrans": reqOrder.ResponseMidtrans,
		"create_time":       reqOrder.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateOrder, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var orderID int64
	err = sc.q.QueryRowx(query, args...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (sc *storeClient) GetOrderByID(ctx context.Context, id int64) (order.Order, error) {
	query := fmt.Sprintf(queryGetOrder, "WHERE t.id = $1")
	// query single row
	var odb orderDB
	err := sc.q.QueryRowx(query, id).StructScan(&odb)
	if err != nil {
		if err == sql.ErrNoRows {
			return order.Order{}, order.ErrDataNotFound
		}
		return order.Order{}, err
	}

	return odb.format(), nil
}
