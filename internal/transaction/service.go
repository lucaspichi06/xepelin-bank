package transaction

import (
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"github.com/lucaspichi06/xepelin-bank/internal/events"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"time"
)

type Service interface {
	Create(tr *domain.Transaction) error
}

type service struct {
	s account.Service
	r Repository
}

func NewService(r Repository, s account.Service) Service {
	return &service{
		r: r,
		s: s,
	}
}

func (s service) Create(tr *domain.Transaction) error {
	var event domain.Event
	switch tr.Type {
	case domain.Deposit:
		tr.DestinationID = nil
		event = events.NewDepositEvent(tr.AccountID, tr.Amount, s.s)
	case domain.WithDraw:
		tr.DestinationID = nil
		event = events.NewWithdrawEvent(tr.AccountID, tr.Amount, s.s)
	case domain.Transfer:
		if tr.DestinationID == nil {
			return custom_errors.ErrInvalidTransactionDestination
		}
		event = events.NewTransferEvent(tr.AccountID, *tr.DestinationID, tr.Amount, s.s)
	default:
		return custom_errors.ErrInvalidTransactionType
	}
	tr.ID = uuid.New()
	tr.Timestamp = time.Now().Format(time.RFC850)

	if _, err := event.Process(); err != nil {
		return err
	}
	return s.r.Create(tr)
}
