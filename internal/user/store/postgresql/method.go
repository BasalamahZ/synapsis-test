package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/synapsis-test/internal/user"
)

// CreateUser insert the given user
//
// CreateUser return created user ID
func (sc *storeClient) CreateUser(ctx context.Context, reqUser user.User) (int64, error) {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"email":        reqUser.Email,
		"name":         reqUser.Name,
		"password":     reqUser.Password,
		"phone_number": reqUser.PhoneNumber,
		"create_time":  reqUser.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateUser, argsKV)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = sc.q.Rebind(query)

	// execute query
	var userID int64
	err = sc.q.QueryRowx(query, args...).Scan(&userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr != nil {
			if pqErr.Code.Name() == "unique_violation" {
				return 0, user.ErrUserAlreadyExist
			}
		}
		return 0, err
	}

	return userID, nil
}

// GetUserByID selects user with the given user ID.
func (sc *storeClient) GetUserByID(ctx context.Context, userID int64) (user.User, error) {
	query := fmt.Sprintf(queryGetUser, "u.id = $1")
	// query single row
	var udb userDB
	err := sc.q.QueryRowx(query, userID).StructScan(&udb)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrDataNotFound
		}
		return user.User{}, err
	}

	return udb.format(), nil
}

// GetUserByEmail select a user with the given
// email.
func (sc *storeClient) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	query := fmt.Sprintf(queryGetUser, "u.email = $1")
	// query single row
	var udb userDB
	err := sc.q.QueryRowx(query, email).StructScan(&udb)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrDataNotFound
		}
		return user.User{}, err
	}

	return udb.format(), nil
}

// UpdateUser updates existing data with the given data
// for a user specified with the given user ID.
//
// UpdateUser do updates on all the fields except ID, Email,
// PICName, CompanyName, PhoneNumber, and CreateTime. So, make
// sure to use current values in the given data if do
// not want to update some specific fields.
func (sc *storeClient) UpdateUser(ctx context.Context, reqUser user.User) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":          reqUser.ID,
		"password":    reqUser.Password,
		"update_time": reqUser.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateUser, argsKV)
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
		if pqErr, ok := err.(*pq.Error); ok && pqErr != nil {
			if pqErr.Code.Name() == "unique_violation" {
				return user.ErrUserAlreadyExist
			}
		}
		return err
	}

	return err
}
