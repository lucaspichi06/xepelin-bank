package account

import (
	"github.com/lucaspichi06/xepelin-bank/internal/domain"

	"github.com/google/uuid"
)

type Service interface {
	Create(account domain.Account) error
	Read(id uuid.UUID) (domain.Account, error)
	Update(account domain.Account) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s service) Create(account domain.Account) error {
	return s.r.Create(account)
}

func (s service) Read(id uuid.UUID) (domain.Account, error) {
	return s.r.Read(id)
}

func (s service) Update(account domain.Account) error {
	return s.r.Update(account)
}
