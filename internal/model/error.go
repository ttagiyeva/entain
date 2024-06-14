package model

import "errors"

const (
	//ErrorBadRequest will throw if the request is not valid
	ErrorBadRequest = "Bad Request"
)

var (
	// ErrorInternalServerError will throw if any the Internal Server Error happen
	ErrorInternalServerError = errors.New("internal server error")
	// ErrorUserNotFound will throw if the requested user is not found
	ErrorUserNotFound = errors.New("user not found")
	// ErrorInsufficientBalance will throw if the request cannot be processed due to insufficient balance
	ErrorInsufficientBalance = errors.New("insufficient balance error")
	// ErrorTransactionAlreadyExists will throw if the given transactionId param has already been processed
	ErrorTransactionAlreadyExists = errors.New("transactionID already exists")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
