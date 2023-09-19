package transaction

import "github.com/lucaspichi06/xepelin-bank/internal/domain"

type Service interface {
	Create(transaction domain.Transaction) (int64, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s service) Create(transaction domain.Transaction) (int64, error) {
	return s.r.Create(transaction)
}
