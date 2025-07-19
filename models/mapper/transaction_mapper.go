package mapper

import (
	"internal-transfer/models"
	"strconv"
)

func CreateTransactionRequestToCreateTransactionInput(req models.CreateTransactionRequest) (*models.TransactionRequest, error) {
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, err
	}

	return &models.TransactionRequest{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               amount,
		IdempotencyKey:       req.IdempotencyKey,
	}, nil
}
