package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	dbCon *sqlx.DB
	tx    *sqlx.Tx
}

func (db Db) NewTransaction(ctx context.Context, options *sql.TxOptions) (*Transaction, error) {
	var err error

	tr := &Transaction{}

	tr.dbCon = db.GetConn()

	tr.tx, err = tr.dbCon.BeginTxx(ctx, options)
	if err != nil {
		if tr.tx != nil {
			_ = tr.tx.Rollback()
		}
		return nil, err
	}

	return tr, nil
}

func (tr *Transaction) PersistNamedCtx(ctx context.Context, query string, entity interface{}) error {
	_, err := tr.tx.NamedExecContext(ctx, query, entity)
	if err != nil {
		return err
	}

	return nil
}

func (tr *Transaction) PersistExecContext(ctx context.Context, query string) error {
	_, err := tr.tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (tr *Transaction) Rollback() error {
	if err := tr.tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (tr *Transaction) Commit() error {
	if err := tr.tx.Commit(); err != nil {
		return err
	}

	return nil
}
