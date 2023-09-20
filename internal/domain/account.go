package domain

import "github.com/google/uuid"

type Account struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance float64   `json:"balance"`
}

type AccountRequest struct {
	Name string `json:"name" binding:"required"`
}
