package transaction

import (
	"database/sql"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
)

type Repository interface {
	Create(transaction domain.Transaction) (int64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(transaction domain.Transaction) (int64, error) {
	query := "INSERT INTO transactions (account_id, type, amount, timestamp) VALUES (?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(transaction.AccountID, transaction.Type, transaction.Amount, transaction.Timestamp)
	if err != nil {
		return 0, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
