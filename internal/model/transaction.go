package model

import "time"

type Transaction struct {
	TransactionID string  `json:"transactionId"`
	State         string  `json:"state"`
	Amount        float32 `json:"amount"`
	UserID        string
	SourceType    string
}

// TransactionDao is the domain object for transactions table.
type TransactionDao struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	TransactionID string    `db:"transaction_id"`
	SourceType    string    `db:"source_type"`
	State         string    `db:"state"`
	Amount        float32   `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
	Cancelled     bool      `db:"cancelled"`
	CancelledAt   time.Time `db:"cancelled_at"`
}
