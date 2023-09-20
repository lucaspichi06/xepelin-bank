package transaction

import (
	"database/sql"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
)

type Repository interface {
	Create(tr *domain.Transaction) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(tr *domain.Transaction) error {

	query := "INSERT INTO transactions (id, account_id, destination_id, type, amount, timestamp) VALUES (?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(tr.ID, tr.AccountID, tr.DestinationID, tr.Type, tr.Amount, tr.Timestamp)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
