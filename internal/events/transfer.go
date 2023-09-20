package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"reflect"
)

type transferEvent struct {
	domain.DefaultEvent
	TargetId uuid.UUID
	Amount   float64
	service  account.Service
}

func NewTransferEvent(id uuid.UUID, targetId uuid.UUID, amt float64, service account.Service) domain.Event {
	var event transferEvent
	event.AccId = id
	event.Type = domain.Transfer
	event.Amount = amt
	event.TargetId = targetId
	event.service = service
	return &event
}

func (t *transferEvent) Process() (domain.Account, error) {
	acc, err := t.service.Read(t.AccId)
	if err != nil {
		return domain.Account{}, err
	}

	if reflect.DeepEqual(acc, domain.Account{}) {
		return domain.Account{}, custom_errors.ErrNotFound
	}

	destAcc, err := t.service.Read(t.TargetId)
	if err != nil {
		return domain.Account{}, err
	}

	if reflect.DeepEqual(destAcc, domain.Account{}) {
		return domain.Account{}, custom_errors.ErrNotFound
	}

	if acc.Balance < t.Amount {
		return domain.Account{}, custom_errors.ErrInsuficientBalance
	}

	acc.Balance -= t.Amount
	destAcc.Balance += t.Amount
	if err = t.service.Update(acc); err != nil {
		return domain.Account{}, err
	}

	if err = t.service.Update(destAcc); err != nil {
		return domain.Account{}, err
	}

	return acc, nil
}
