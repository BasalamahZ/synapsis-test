package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/category"
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

// categoryDB denotes a data in the store.
type categoryDB struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreateTime  time.Time  `db:"create_time"`
	UpdateTime  *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (cdb *categoryDB) format() category.Category {
	p := category.Category{
		ID:          cdb.ID,
		Name:        cdb.Name,
		Description: cdb.Description,
		CreateTime:  cdb.CreateTime,
	}

	if cdb.UpdateTime != nil {
		p.UpdateTime = *cdb.UpdateTime
	}

	return p
}
