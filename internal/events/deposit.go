package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"reflect"
)

type depositEvent struct {
	domain.DefaultEvent
	Amount  float64
	service account.Service
}

func NewDepositEvent(id uuid.UUID, amt float64, service account.Service) domain.Event {
	var event depositEvent
	event.AccId = id
	event.Type = domain.Deposit
	event.Amount = amt
	event.service = service
	return &event
}

func (t *depositEvent) Process() (domain.Account, error) {
	acc, err := t.service.Read(t.AccId)
	if err != nil {
		return domain.Account{}, err
	}

	if reflect.DeepEqual(acc, domain.Account{}) {
		return domain.Account{}, custom_errors.ErrNotFound
	}

	acc.Balance = acc.Balance + t.Amount
	return acc, t.service.Update(acc)
}
