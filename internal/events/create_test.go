package events

import (
	"errors"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProcess(t *testing.T) {
	t.Run("create process success", func(t *testing.T) {
		serviceMock := accServiceMock{
			create: func(account domain.Account) error {
				return nil
			},
		}

		create := NewCreateAccountEvent("test", serviceMock)

		acc, err := create.Process()

		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, "test", acc.Name)
		assert.Equal(t, 0.00, acc.Balance)
	})
	t.Run("balance process error", func(t *testing.T) {
		serviceMock := accServiceMock{
			create: func(account domain.Account) error {
				return errors.New("test error")
			},
		}

		create := NewCreateAccountEvent("test", serviceMock)

		_, err := create.Process()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
}
