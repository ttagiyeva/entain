package model

// TransactionToTransactionDao converts a transaction to a transaction dao.
func TransactionToTransactionDao(t *Transaction) *TransactionDao {
	return &TransactionDao{
		UserID:        t.UserID,
		TransactionID: t.TransactionID,
		SourceType:    t.SourceType,
		State:         t.State,
		Amount:        t.Amount,
	}
}
