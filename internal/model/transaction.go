package model

type Transaction struct {
	ID            string
	TransactionID string  `json:"transactionId"`
	State         string  `json:"state"`
	Amount        float32 `json:"amount"`
	UserID        string
	SourceType    string
}

var SourceType = map[string]struct{}{"game": {}, "server": {}, "payment": {}}

var State = map[string]string{"win": Win, "lost": Lost}

const (
	Win  = "win"
	Lost = "lost"
)
