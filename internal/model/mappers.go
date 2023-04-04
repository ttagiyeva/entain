package model

// TransactionDaoToTransaction converts a transaction dao to a transaction.
func TransactionDaoToTransaction(t *TransactionDao) *Transaction {
	return &Transaction{
		UserID:        t.UserID,
		TransactionID: t.TransactionID,
		SourceType:    t.SourceType,
		State:         t.State,
		Amount:        t.Amount,
	}
}

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
