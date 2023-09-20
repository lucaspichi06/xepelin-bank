package events

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
)

type createEvent struct {
	domain.DefaultEvent
	AccName string
	service account.Service
}

func NewCreateAccountEvent(name string, service account.Service) domain.Event {
	var event createEvent
	event.AccId = uuid.New()
	event.Type = domain.Create
	event.AccName = name
	event.service = service
	return &event
}

func (t *createEvent) Process() (domain.Account, error) {
	acc := domain.Account{
		ID:      t.AccId,
		Name:    t.AccName,
		Balance: 0,
	}
	return acc, t.service.Create(acc)
}
