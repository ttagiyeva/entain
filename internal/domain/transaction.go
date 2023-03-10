package domain

import "time"

// Transaction is the domain object for transactions table.
type Transaction struct {
	ID            string     `db:"id"`
	UserID        string     `db:"user_id"`
	TransactionID string     `db:"transaction_id"`
	SourceType    string     `db:"source_type"`
	State         string     `db:"state"`
	Amount        float32    `db:"amount"`
	CreatedAt     *time.Time `db:"created_at"`
	Cancelled     bool       `db:"cancelled"`
	CancelledAt   *time.Time `db:"cancelled_at"`
}
