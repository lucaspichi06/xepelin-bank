package events

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransferProcess(t *testing.T) {
	t.Run("transfer process success", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), uuid.New(), 100.00, serviceMock)

		acc, err := transfer.Process()

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, 900.00, acc.Balance)
	})
	t.Run("transfer process not found", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, custom_errors.ErrNotFound
			},
		}

		transfer := NewTransferEvent(uuid.New(), uuid.New(), 100.00, serviceMock)

		acc, err := transfer.Process()

		assert.Error(t, err)
		assert.Equal(t, err, custom_errors.ErrNotFound)
		assert.Equal(t, domain.Account{}, acc)
	})
	t.Run("transfer process empty account", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), uuid.New(), 100.00, serviceMock)

		acc, err := transfer.Process()

		assert.Error(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, domain.Account{}, acc)
		assert.Equal(t, err, custom_errors.ErrNotFound)
	})
	t.Run("transfer process not found destination account", func(t *testing.T) {
		destination := uuid.New()
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				if id == destination {
					return domain.Account{}, custom_errors.ErrNotFound
				}
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), destination, 100.00, serviceMock)

		acc, err := transfer.Process()

		assert.Error(t, err)
		assert.Equal(t, err, custom_errors.ErrNotFound)
		assert.Equal(t, domain.Account{}, acc)
	})
	t.Run("transfer process empty destination account", func(t *testing.T) {
		destination := uuid.New()
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				if id == destination {
					return domain.Account{}, nil
				}
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), destination, 100.00, serviceMock)

		acc, err := transfer.Process()

		assert.Error(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, domain.Account{}, acc)
		assert.Equal(t, err, custom_errors.ErrNotFound)
	})
	t.Run("transfer process negative balance", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 0.00,
				}, nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), uuid.New(), 100.00, serviceMock)

		_, err := transfer.Process()

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrInsuficientBalance, err)
	})
	t.Run("transfer process update error", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
			update: func(account domain.Account) error {
				return errors.New("test error")
			},
		}

		transfer := NewTransferEvent(uuid.New(), uuid.New(), 100.00, serviceMock)

		_, err := transfer.Process()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
	t.Run("transfer process update error destination account", func(t *testing.T) {
		destination := uuid.New()
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
			update: func(account domain.Account) error {
				if account.ID == destination {
					return errors.New("test error")
				}
				return nil
			},
		}

		transfer := NewTransferEvent(uuid.New(), destination, 100.00, serviceMock)

		_, err := transfer.Process()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
}
