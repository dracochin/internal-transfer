package mapper

import (
	"internal-transfer/models"
	"strconv"
)

func CreateAccountRequestToCreateAccountInput(req models.CreateAccountRequest) (*models.Account, error) {
	balance, err := strconv.ParseFloat(req.Balance, 64)
	if err != nil {
		return nil, err
	}

	return &models.Account{
		AccountID: req.AccountID,
		Balance:   balance,
	}, nil
}
