package transaction

import (
	_ "database/sql"
	"errors"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("create transaction success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		tr := domain.Transaction{
			ID:        uuid.New(),
			AccountID: uuid.New(),
			Type:      domain.Deposit,
		}

		err = repo.Create(&tr)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("create transaction prepare error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO transactions").
			WillReturnError(errors.New("test error"))

		tr := domain.Transaction{
			ID:        uuid.New(),
			AccountID: uuid.New(),
			Type:      domain.Deposit,
		}

		err = repo.Create(&tr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("create transaction exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fail()
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("test error"))

		tr := domain.Transaction{
			ID:        uuid.New(),
			AccountID: uuid.New(),
			Type:      domain.Deposit,
		}

		err = repo.Create(&tr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
