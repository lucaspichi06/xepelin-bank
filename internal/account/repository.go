package account

import (
	"database/sql"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
)

type Repository interface {
	Create(account domain.Account) (int64, error)
	Read(id int) (domain.Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(account domain.Account) (int64, error) {
	query := "INSERT INTO accounts (name, balance) VALUES (?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(account.Name, account.Balance)
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

func (r repository) Read(id int) (domain.Account, error) {
	var account domain.Account
	query := "SELECT * FROM accounts WHERE id = ?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&account.ID, &account.Name, &account.Balance)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}
