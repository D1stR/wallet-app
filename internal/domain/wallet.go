package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type OperationType int

const (
	Withdraw OperationType = iota
	Deposit
)

func (o OperationType) String() string {
	switch o {
	case Withdraw:
		return "WITHDRAW"
	case Deposit:
		return "DEPOSIT"
	default:
		return "UNKNOWN"
	}
}

func OperationTypeFromString(opType string) (OperationType, error) {
	switch opType {
	case "DEPOSIT":
		return Deposit, nil
	case "WITHDRAW":
		return Withdraw, nil
	default:
		return -1, fmt.Errorf("invalid operation type: %s", opType)
	}
}

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}
