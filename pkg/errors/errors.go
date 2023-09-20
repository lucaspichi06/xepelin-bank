package custom_errors

import "errors"

var (
	// handler errors
	ErrNotFound    = errors.New("registry not found")
	ErrInvalidJSON = errors.New("invalid json")
	ErrInvalidID   = errors.New("invalid param id")

	// account errors
	ErrAccountExist       = errors.New("there is already an account with this name")
	ErrInsuficientBalance = errors.New("insufficient amount in the account balance")

	// transaction errors
	ErrInvalidTransactionType        = errors.New("invalid transaction type")
	ErrInvalidTransactionDestination = errors.New("invalid transaction destination account")
)
