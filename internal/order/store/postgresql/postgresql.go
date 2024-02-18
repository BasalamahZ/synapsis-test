package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/order"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements user/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements user/service.PGStoreClient
type storeClient struct {
	q sqlx.Ext
}

// New creates a new store.
func New(db *sqlx.DB) (*store, error) {
	s := &store{
		db: db,
	}

	return s, nil
}

func (s *store) NewClient(useTx bool) (PGStoreClient, error) {
	var q sqlx.Ext

	// determine what object should be use as querier
	q = s.db
	if useTx {
		var err error
		q, err = s.db.Beginx()
		if err != nil {
			return nil, err
		}
	}

	return &storeClient{
		q: q,
	}, nil
}

func (sc *storeClient) Commit() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Commit()
	}
	return errInvalidCommit
}

func (sc *storeClient) Rollback() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Rollback()
	}
	return errInvalidRollback
}

// orderDB denotes a data in the store.
type orderDB struct {
	ID               int64        `db:"id"`
	UserID           int64        `db:"user_id"`
	UserName         string       `db:"user_name"`
	UserEmail        string       `db:"user_email"`
	ProductID        int64        `db:"product_id"`
	ProductName      string       `db:"product_name"`
	ProductPrice     int64        `db:"product_price"`
	Quantity         int64        `db:"quantity"`
	TotalAmount      int64        `db:"total_amount"`
	Status           order.Status `db:"status"`
	ResponseMidtrans string       `db:"response_midtrans"`
	CreateTime       time.Time    `db:"create_time"`
	UpdateTime       *time.Time   `db:"update_time"`
}

// format formats database struct into domain struct.
func (odb *orderDB) format() order.Order {
	o := order.Order{
		ID:               odb.ID,
		UserID:           odb.UserID,
		UserName:         odb.UserName,
		UserEmail:        odb.UserEmail,
		ProductID:        odb.ProductID,
		ProductName:      odb.ProductName,
		ProductPrice:     odb.ProductPrice,
		Quantity:         odb.Quantity,
		TotalAmount:      odb.TotalAmount,
		Status:           odb.Status,
		ResponseMidtrans: odb.ResponseMidtrans,
		CreateTime:       odb.CreateTime,
	}

	if odb.UpdateTime != nil {
		o.UpdateTime = *odb.UpdateTime
	}

	return o
}
