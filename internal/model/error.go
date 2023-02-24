package model

import "errors"

const (
	//ErrorInvalidSourceType will throw if the given request header source type is not valid
	ErrorInvalidSourceType = "Source type is invalid"
	//ErrorInvalidSourceType will throw if the given request header state is not valid
	ErrorInvalidState = "State is invalid"
	// ErrorInvalidAmount will throw if the given request body amount is not valid
	ErrorInvalidAmont = "Amount is invalid"
	//ErrorInvalidTransactionId will throw if the given request body transactionId is not valid
	ErrorInvalidTransactionId = "TransactionId is invalid"
	//ErrorBadRequest will throw if the given request param is not valid
	ErrorBadRequest = "Bad Request"
)

var (
	// ErrorInternalServer will throw if any the Internal Server Error has happen
	ErrorInternalServer = errors.New("internal Server Error")
	// ErrorNotFound will throw if the requested user does not exist
	ErrorNotFound = errors.New("requested user is not found")
	// ErrorInsufficientBalance will throw if the request cannot be processed due to insufficient balance
	ErrorInsufficientBalance = errors.New("insufficient balance")
	// ErrorTransactionAlreadyExists will throw if the given request param is not valid
	ErrorTransactionAlreadyExists = errors.New("transactionID already exists")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
