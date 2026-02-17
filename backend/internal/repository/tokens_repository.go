package repository

import (
	"database/sql"
	"fmt"

	"github.com/smartstocks/backend/internal/models"
)

type TokensRepository struct {
	db *sql.DB
}

func NewTokensRepository(db *sql.DB) *TokensRepository {
	return &TokensRepository{db: db}
}

// GetUserTokens obtiene el balance de tokens del usuario
func (r *TokensRepository) GetUserTokens(userID string) (*models.UserTokens, error) {
	tokens := &models.UserTokens{}

	query := `
		SELECT user_id, balance, total_earned, total_spent, 
			   last_transaction_at, created_at, updated_at
		FROM user_tokens
		WHERE user_id = ?
	`

	err := r.db.QueryRow(query, userID).Scan(
		&tokens.UserID,
		&tokens.Balance,
		&tokens.TotalEarned,
		&tokens.TotalSpent,
		&tokens.LastTransactionAt,
		&tokens.CreatedAt,
		&tokens.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user tokens not found")
	}

	return tokens, err
}

// AddTokens añade tokens al usuario
func (r *TokensRepository) AddTokens(userID string, amount int, transactionType, description string, referenceID *string) error {
	var refID sql.NullString
	if referenceID != nil {
		refID = sql.NullString{String: *referenceID, Valid: true}
	}

	query := `CALL add_tokens(?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, userID, amount, transactionType, description, refID)
	return err
}

// SubtractTokens resta tokens al usuario
func (r *TokensRepository) SubtractTokens(userID string, amount int, transactionType, description string, referenceID *string) (bool, error) {
	var refID sql.NullString
	if referenceID != nil {
		refID = sql.NullString{String: *referenceID, Valid: true}
	}

	var success bool
	query := `CALL subtract_tokens(?, ?, ?, ?, ?, @success)`
	_, err := r.db.Exec(query, userID, amount, transactionType, description, refID)
	if err != nil {
		return false, err
	}

	err = r.db.QueryRow(`SELECT @success`).Scan(&success)
	return success, err
}

// GetTransactionHistory obtiene el historial de transacciones
func (r *TokensRepository) GetTransactionHistory(userID string, limit int) ([]models.TokenTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
		SELECT id, user_id, transaction_type, amount, balance_after,
			   description, reference_id, created_at
		FROM token_transactions
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.TokenTransaction
	for rows.Next() {
		var tx models.TokenTransaction
		var refID sql.NullString

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.TransactionType,
			&tx.Amount,
			&tx.BalanceAfter,
			&tx.Description,
			&refID,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if refID.Valid {
			tx.ReferenceID = &refID.String
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// HasSufficientTokens verifica si el usuario tiene suficientes tokens
func (r *TokensRepository) HasSufficientTokens(userID string, amount int) (bool, error) {
	var balance int
	query := `SELECT balance FROM user_tokens WHERE user_id = ?`
	err := r.db.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		return false, err
	}
	return balance >= amount, nil
}
