package events

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithdrawProcess(t *testing.T) {
	t.Run("withdraw process success", func(t *testing.T) {
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

		withdraw := NewWithdrawEvent(uuid.New(), 100.00, serviceMock)

		acc, err := withdraw.Process()

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, 900.00, acc.Balance)
	})
	t.Run("withdraw process not found", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, custom_errors.ErrNotFound
			},
		}

		withdraw := NewWithdrawEvent(uuid.New(), 100.00, serviceMock)

		acc, err := withdraw.Process()

		assert.Error(t, err)
		assert.Equal(t, err, custom_errors.ErrNotFound)
		assert.Equal(t, domain.Account{}, acc)
	})
	t.Run("withdraw process empty account", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, nil
			},
		}

		withdraw := NewWithdrawEvent(uuid.New(), 100.00, serviceMock)

		acc, err := withdraw.Process()

		assert.Error(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, domain.Account{}, acc)
		assert.Equal(t, err, custom_errors.ErrNotFound)
	})
	t.Run("withdraw process negative balance", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 0.00,
				}, nil
			},
		}

		withdraw := NewWithdrawEvent(uuid.New(), 100.00, serviceMock)

		_, err := withdraw.Process()

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrInsuficientBalance, err)
	})
	t.Run("withdraw process update error", func(t *testing.T) {
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

		withdraw := NewWithdrawEvent(uuid.New(), 100.00, serviceMock)

		_, err := withdraw.Process()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
}
