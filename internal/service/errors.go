package service

import "errors"

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInvalidOperationType = errors.New("invalid operation type")
)
