package model

import "time"

// CreateTransaction contain fields to create a transaction entry
type CreateTransaction struct {
	UserId      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// Transaction contain fields of transactions. The fields refers to 'transactions' table
type Transaction struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Status      string    `json:"status" db:"status"`
	Description string    `json:"description" db:"description"`
	Remarks     string    `json:"remarks" db:"remarks"`
	Amount      float64   `json:"amount" db:"amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
