package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/product"
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

// productDB denotes a data in the store.
type productDB struct {
	ID           int64      `db:"id"`
	Name         string     `db:"name"`
	Price        int64      `db:"price"`
	Description  string     `db:"description"`
	CategoryID   int64      `db:"category_id"`
	CategoryName string     `db:"category_name"`
	CreateTime   time.Time  `db:"create_time"`
	UpdateTime   *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (pdb *productDB) format() product.Product {
	p := product.Product{
		ID:           pdb.ID,
		Name:         pdb.Name,
		Price:        pdb.Price,
		Description:  pdb.Description,
		CategoryID:   pdb.CategoryID,
		CategoryName: pdb.CategoryName,
		CreateTime:   pdb.CreateTime,
	}

	if pdb.UpdateTime != nil {
		p.UpdateTime = *pdb.UpdateTime
	}

	return p
}

// productCartDB denotes a data in the store.
type productCartDB struct {
	UserID       int64      `db:"user_id"`
	ProductID    int64      `db:"product_id"`
	ProductName  string     `db:"product_name"`
	ProductPrice int64      `db:"product_price"`
	Quantity     int64      `db:"quantity"`
	CreateTime   time.Time  `db:"create_time"`
	UpdateTime   *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (pcdb *productCartDB) format() product.ProductCart {
	pc := product.ProductCart{
		UserID:       pcdb.UserID,
		ProductID:    pcdb.ProductID,
		ProductName:  pcdb.ProductName,
		ProductPrice: pcdb.ProductPrice,
		Quantity:     pcdb.Quantity,
		CreateTime:   pcdb.CreateTime,
	}

	if pcdb.UpdateTime != nil {
		pc.UpdateTime = *pcdb.UpdateTime
	}

	return pc
}
