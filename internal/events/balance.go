package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"reflect"
)

type balanceEvent struct {
	domain.DefaultEvent
	service account.Service
}

func NewBalanceEvent(id uuid.UUID, service account.Service) domain.Event {
	var event balanceEvent
	event.AccId = id
	event.Type = domain.Balance
	event.service = service
	return &event
}

func (t *balanceEvent) Process() (domain.Account, error) {
	acc, err := t.service.Read(t.AccId)
	if err != nil {
		return domain.Account{}, err
	}

	if reflect.DeepEqual(acc, domain.Account{}) {
		return domain.Account{}, custom_errors.ErrNotFound
	}

	return acc, nil
}
