package postgresql

import (
	"errors"
	"time"

	"github.com/synapsis-test/internal/user"
	"github.com/jmoiron/sqlx"
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

// userDB denotes a data in the store.
type userDB struct {
	ID          int64      `db:"id"`
	Email       string     `db:"email"`
	Name        string     `db:"name"`
	Password    string     `db:"password"`
	PhoneNumber string     `db:"phone_number"`
	CreateTime  time.Time  `db:"create_time"`
	UpdateTime  *time.Time `db:"update_time"`
}

// format formats database struct into domain struct.
func (udb *userDB) format() user.User {
	u := user.User{
		ID:          udb.ID,
		Email:       udb.Email,
		Name:        udb.Name,
		Password:    udb.Password,
		PhoneNumber: udb.PhoneNumber,
		CreateTime:  udb.CreateTime,
	}

	if udb.UpdateTime != nil {
		u.UpdateTime = *udb.UpdateTime
	}

	return u
}
