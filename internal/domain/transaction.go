package domain

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	AccountID string    `json:"account_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
