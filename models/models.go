package models

type Account struct {
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type TransactionRequest struct {
	SourceAccountID      int64   `json:"source_account_id"`
	DestinationAccountID int64   `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
	IdempotencyKey       string  `json:"idempotency_key"`
}
