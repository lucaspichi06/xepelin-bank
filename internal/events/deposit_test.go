package events

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepositProcess(t *testing.T) {
	t.Run("deposit process success", func(t *testing.T) {
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

		deposit := NewDepositEvent(uuid.New(), 100.00, serviceMock)

		acc, err := deposit.Process()

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, 1100.00, acc.Balance)
	})
	t.Run("deposit process not found", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, custom_errors.ErrNotFound
			},
		}

		deposit := NewDepositEvent(uuid.New(), 100.00, serviceMock)

		acc, err := deposit.Process()

		assert.Error(t, err)
		assert.Equal(t, err, custom_errors.ErrNotFound)
		assert.Equal(t, domain.Account{}, acc)
	})
	t.Run("deposit process empty account", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, nil
			},
		}

		deposit := NewDepositEvent(uuid.New(), 100.00, serviceMock)

		acc, err := deposit.Process()

		assert.Error(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, domain.Account{}, acc)
		assert.Equal(t, err, custom_errors.ErrNotFound)
	})
	t.Run("deposit process update error", func(t *testing.T) {
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

		deposit := NewDepositEvent(uuid.New(), 100.00, serviceMock)

		_, err := deposit.Process()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
}
