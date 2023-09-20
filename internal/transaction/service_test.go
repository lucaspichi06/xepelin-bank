package transaction

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

type trRepositoryMock struct {
	create func(tr *domain.Transaction) error
}

func (t trRepositoryMock) Create(tr *domain.Transaction) error {
	return t.create(tr)
}

func TestTransactionCreate(t *testing.T) {
	t.Run("transaction deposit success", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID: uuid.New(),
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		tr := domain.Transaction{
			Type: domain.Deposit,
		}

		err := trService.Create(&tr)

		assert.NoError(t, err)
	})
	t.Run("transaction withdraw success", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID: uuid.New(),
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		tr := domain.Transaction{
			Type: domain.WithDraw,
		}

		err := trService.Create(&tr)

		assert.NoError(t, err)
	})
	t.Run("transaction transfer success", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID: uuid.New(),
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		id := uuid.New()
		tr := domain.Transaction{
			Type:          domain.Transfer,
			DestinationID: &id,
		}

		err := trService.Create(&tr)

		assert.NoError(t, err)
	})
	t.Run("transaction transfer error - invalid destination id", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID: uuid.New(),
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		tr := domain.Transaction{
			Type: domain.Transfer,
		}

		err := trService.Create(&tr)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrInvalidTransactionDestination, err)
	})
	t.Run("transaction transfer error - invalid account id", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		id := uuid.New()
		tr := domain.Transaction{
			Type:          domain.Transfer,
			DestinationID: &id,
		}

		err := trService.Create(&tr)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
	t.Run("transaction invalid type error", func(t *testing.T) {
		serviceMock := accServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					ID: uuid.New(),
				}, nil
			},
			update: func(account domain.Account) error {
				return nil
			},
		}

		repoMock := trRepositoryMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}

		trService := NewService(repoMock, serviceMock)
		tr := domain.Transaction{
			Type: domain.Create,
		}

		err := trService.Create(&tr)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrInvalidTransactionType, err)
	})
}
