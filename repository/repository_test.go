package repository_test

import (
	"regexp"
	"testing"

	"internal-transfer/models"
	"internal-transfer/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	account := models.Account{
		AccountID: 123,
		Balance:   100.00,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`)).
		WithArgs(account.AccountID, account.Balance).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateAccount(account)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	rows := sqlmock.NewRows([]string{"balance"}).
		AddRow("100.00")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT balance FROM accounts WHERE account_id = $1`)).
		WithArgs(123).
		WillReturnRows(rows)

	account, err := repo.GetAccount(123)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), account.AccountID)
	assert.Equal(t, 100.00, account.Balance)
}
