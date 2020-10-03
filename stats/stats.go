package stats

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

// Database represents a database.
type Database struct {
	*sqlx.DB
}

func Connect(addr string) (*Database, error) {
	d, err := sqlx.Open("mysql", addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}

	d.SetConnMaxLifetime(5 * time.Minute)
	d.SetMaxOpenConns(100)
	d.SetMaxIdleConns(100)

	return &Database{d}, nil
}

// This used to be a transaction.

type ReadView struct {
	*Database
	ctx context.Context
}

var roOpts = &sql.TxOptions{
	ReadOnly: true,
}

func (db *Database) WithContext(ctx context.Context) ReadView {
	return ReadView{db, ctx}
}
