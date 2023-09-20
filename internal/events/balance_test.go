package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type accServiceMock struct {
	create func(account domain.Account) error
	read   func(id uuid.UUID) (domain.Account, error)
	update func(account domain.Account) error
}

func (a accServiceMock) Create(account domain.Account) error {
	return a.create(account)
}

func (a accServiceMock) Read(id uuid.UUID) (domain.Account, error) {
	return a.read(id)
}

func (a accServiceMock) Update(account domain.Account) error {
	return a.update(account)
}

func TestBalanceProcess(t *testing.T) {
	t.Run("balance process success", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID:      id,
					Balance: 1000.00,
				}, nil
			},
		}

		balance := NewBalanceEvent(uuid.New(), serviceMock)

		acc, err := balance.Process()

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, 1000.00, acc.Balance)
	})
	t.Run("balance process not found", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, custom_errors.ErrNotFound
			},
		}

		balance := NewBalanceEvent(uuid.New(), serviceMock)

		acc, err := balance.Process()

		assert.Error(t, err)
		assert.Equal(t, err, custom_errors.ErrNotFound)
		assert.Equal(t, domain.Account{}, acc)
	})
	t.Run("balance process empty account", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, nil
			},
		}

		balance := NewBalanceEvent(uuid.New(), serviceMock)

		acc, err := balance.Process()

		assert.Error(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, domain.Account{}, acc)
		assert.Equal(t, err, custom_errors.ErrNotFound)
	})
}
