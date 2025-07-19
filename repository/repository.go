package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"internal-transfer/models"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateAccount(acc models.Account) error {
	_, err := r.DB.Exec("INSERT INTO accounts (account_id, balance) VALUES ($1, $2)", acc.AccountID, acc.Balance)
	return err
}

func (r *Repository) GetAccount(id int64) (*models.Account, error) {
	var balance float64
	err := r.DB.QueryRow("SELECT balance FROM accounts WHERE account_id = $1", id).Scan(&balance)
	if err == sql.ErrNoRows {
		return nil, errors.New("account not found")
	}
	if err != nil {
		return nil, err
	}
	return &models.Account{AccountID: id, Balance: balance}, nil
}

// Checks if a transaction with this idempotency key exists in the given transaction
func (r *Repository) TransactionExists(tx *sql.Tx, idempotencyKey string) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM transactions WHERE idempotency_key = $1)", idempotencyKey).Scan(&exists)
	return exists, err
}

// Create a transaction and double-entry records atomically inside the given transaction
func (r *Repository) CreateTransactionWithEntries(tx *sql.Tx, t models.TransactionRequest) error {
	if t.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Lock both accounts for update, ordered by account_id to avoid deadlocks
	ids := []int64{t.SourceAccountID, t.DestinationAccountID}
	if ids[0] > ids[1] {
		ids[0], ids[1] = ids[1], ids[0]
	}

	rows, err := tx.Query(`SELECT account_id, balance FROM accounts WHERE account_id IN ($1, $2) FOR UPDATE`, ids[0], ids[1])
	if err != nil {
		return err
	}
	defer rows.Close()

	balances := make(map[int64]float64)
	for rows.Next() {
		var id int64
		var bal float64
		if err := rows.Scan(&id, &bal); err != nil {
			return err
		}
		balances[id] = bal
	}

	if len(balances) < 2 {
		return errors.New("one or both accounts not found")
	}

	if balances[t.SourceAccountID] < t.Amount {
		return errors.New("insufficient funds")
	}

	// Create transaction record
	txnID := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO transactions (id, source_account_id, destination_account_id, amount, idempotency_key)
		VALUES ($1, $2, $3, $4, $5)`,
		txnID, t.SourceAccountID, t.DestinationAccountID, t.Amount, t.IdempotencyKey,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Create entries and update balances
	entries := []struct {
		AccountID int64
		Amount    float64
	}{
		{t.SourceAccountID, -t.Amount},
		{t.DestinationAccountID, t.Amount},
	}

	for _, e := range entries {
		_, err := tx.Exec(
			`INSERT INTO entries (transaction_id, account_id, amount) VALUES ($1, $2, $3)`,
			txnID, e.AccountID, e.Amount,
		)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %w", err)
		}

		_, err = tx.Exec(
			`UPDATE accounts SET balance = balance + $1 WHERE account_id = $2`,
			e.Amount, e.AccountID,
		)
		if err != nil {
			return fmt.Errorf("failed to update account balance: %w", err)
		}
	}

	return nil
}
