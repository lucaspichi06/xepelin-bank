package account

import (
	_ "database/sql"
	"errors"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	t.Run("create account success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO accounts").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Create(account)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("create account prepare error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO accounts").
			WillReturnError(errors.New("test error"))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Create(account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("create account exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO accounts").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("test error"))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Create(account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestReadAccount(t *testing.T) {
	t.Run("read account success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectQuery("SELECT \\* FROM accounts WHERE id = \\?").WithArgs(
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "balance"}).AddRow(
			"123e4567-e89b-12d3-a456-426614174000", "test", 100.0,
		))

		account, err := repo.Read(uuid.New())
		assert.NoError(t, err)
		assert.Equal(t, "test", account.Name)
		assert.Equal(t, 100.0, account.Balance)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("read account scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectQuery("SELECT \\* FROM accounts WHERE id = \\?").WithArgs(
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("test error"))

		account, err := repo.Read(uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
		assert.Equal(t, domain.Account{}, account)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUpdateAccount(t *testing.T) {
	t.Run("update account success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("UPDATE accounts").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Update(account)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("update account prepare error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("UPDATE accounts").
			WillReturnError(errors.New("test error"))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Update(account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("update account exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("UPDATE accounts").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("test error"))

		account := domain.Account{
			ID:      uuid.New(),
			Name:    "test",
			Balance: 100.0,
		}

		err = repo.Update(account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
