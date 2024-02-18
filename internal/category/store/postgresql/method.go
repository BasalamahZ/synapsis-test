package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/category"
)

func (sc *storeClient) GetCategoryByID(ctx context.Context, id int64) (category.Category, error) {
	query := fmt.Sprintf(queryGetCategory, "WHERE c.id = $1")
	// query single row
	var cdb categoryDB
	err := sc.q.QueryRowx(query, id).StructScan(&cdb)
	if err != nil {
		if err == sql.ErrNoRows {
			return category.Category{}, category.ErrDataNotFound
		}
		return category.Category{}, err
	}

	return cdb.format(), nil
}

func (sc *storeClient) GetCategories(ctx context.Context) ([]category.Category, error) {
	// construct query
	query := fmt.Sprintf(queryGetCategory, "")

	// prepare query
	query, args, err := sqlx.Named(query, map[string]interface{}{})
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

	// read categories
	categories := make([]category.Category, 0)
	for rows.Next() {
		var row categoryDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		categories = append(categories, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
