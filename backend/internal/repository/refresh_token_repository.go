package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(token *models.RefreshToken) error {
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens WHERE token = ? AND expires_at > ?
	`

	err := r.db.QueryRow(query, token, time.Now()).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid or expired refresh token")
	}

	return rt, err
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *RefreshTokenRepository) DeleteUserRefreshTokens(userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *RefreshTokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < ?`
	_, err := r.db.Exec(query, time.Now())
	return err
}
