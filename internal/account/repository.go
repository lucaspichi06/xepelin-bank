package account

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
)

type Repository interface {
	Create(account domain.Account) error
	Read(id uuid.UUID) (domain.Account, error)
	Update(account domain.Account) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(account domain.Account) error {
	query := "INSERT INTO accounts (id, name, balance) VALUES (?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(account.ID, account.Name, account.Balance)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r repository) Read(id uuid.UUID) (domain.Account, error) {
	var account domain.Account
	query := "SELECT * FROM accounts WHERE id = ?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&account.ID, &account.Name, &account.Balance)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func (r repository) Update(account domain.Account) error {
	query := "UPDATE accounts SET name = ?, balance = ? WHERE id = ?;"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(account.Name, account.Balance, account.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
