package account

import "github.com/lucaspichi06/xepelin-bank/internal/domain"

type Service interface {
	Create(account domain.Account) (int64, error)
	Read(id int) (domain.Account, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s service) Create(account domain.Account) (int64, error) {
	return s.r.Create(account)
}

func (s service) Read(id int) (domain.Account, error) {
	return s.r.Read(id)
}
