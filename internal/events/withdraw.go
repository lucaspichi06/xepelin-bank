package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"reflect"
)

type withdrawEvent struct {
	domain.DefaultEvent
	Amount  float64
	service account.Service
}

func NewWithdrawEvent(id uuid.UUID, amt float64, service account.Service) domain.Event {
	var event withdrawEvent
	event.AccId = id
	event.Type = domain.WithDraw
	event.Amount = amt
	event.service = service
	return &event
}

func (t *withdrawEvent) Process() (domain.Account, error) {
	acc, err := t.service.Read(t.AccId)
	if err != nil {
		return domain.Account{}, err
	}

	if reflect.DeepEqual(acc, domain.Account{}) {
		return domain.Account{}, custom_errors.ErrNotFound
	}

	if acc.Balance < t.Amount {
		return domain.Account{}, custom_errors.ErrInsuficientBalance
	}

	acc.Balance = acc.Balance - t.Amount
	return acc, t.service.Update(acc)
}
