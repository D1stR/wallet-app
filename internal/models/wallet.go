package models

type WalletRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}
