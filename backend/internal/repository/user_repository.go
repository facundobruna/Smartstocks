package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, email, password_hash, profile_picture_url, school_id, 
						  email_verified, verification_token)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ProfilePictureURL,
		user.SchoolID,
		user.EmailVerified,
		user.VerificationToken,
	)

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, profile_picture_url, school_id,
			   created_at, updated_at, last_login, email_verified, verification_token
		FROM users WHERE email = ?
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePictureURL,
		&user.SchoolID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.EmailVerified,
		&user.VerificationToken,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, profile_picture_url, school_id,
			   created_at, updated_at, last_login, email_verified
		FROM users WHERE id = ?
	`

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePictureURL,
		&user.SchoolID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.EmailVerified,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id FROM users WHERE username = ?`

	err := r.db.QueryRow(query, username).Scan(&user.ID)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	query := `UPDATE users SET last_login = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *UserRepository) VerifyEmail(token string) error {
	query := `UPDATE users SET email_verified = TRUE, verification_token = NULL WHERE verification_token = ?`
	result, err := r.db.Exec(query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("invalid verification token")
	}

	return nil
}

func (r *UserRepository) UpdateProfile(userID string, req *models.UpdateProfileRequest) error {
	query := `UPDATE users SET `
	args := []interface{}{}
	updates := []string{}

	if req.Username != nil {
		updates = append(updates, "username = ?")
		args = append(args, *req.Username)
	}
	if req.ProfilePictureURL != nil {
		updates = append(updates, "profile_picture_url = ?")
		args = append(args, *req.ProfilePictureURL)
	}
	if req.SchoolID != nil {
		updates = append(updates, "school_id = ?")
		args = append(args, *req.SchoolID)
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += ", updated_at = ? WHERE id = ?"
	args = append(args, time.Now(), userID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *UserRepository) SetPasswordResetToken(email, token string, expires time.Time) error {
	query := `UPDATE users SET reset_token = ?, reset_token_expires = ? WHERE email = ?`
	_, err := r.db.Exec(query, token, expires, email)
	return err
}

func (r *UserRepository) ResetPassword(token, newPasswordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = ?, reset_token = NULL, reset_token_expires = NULL 
		WHERE reset_token = ? AND reset_token_expires > ?
	`
	result, err := r.db.Exec(query, newPasswordHash, token, time.Now())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("invalid or expired reset token")
	}

	return nil
}

func (r *UserRepository) GetUserStats(userID string) (*models.UserStats, error) {
	stats := &models.UserStats{}
	query := `
		SELECT user_id, smartpoints, rank_tier, total_quizzes_completed,
			   total_simulator_games, win_streak, total_wins, total_losses, updated_at
		FROM user_stats WHERE user_id = ?
	`

	err := r.db.QueryRow(query, userID).Scan(
		&stats.UserID,
		&stats.Smartpoints,
		&stats.RankTier,
		&stats.TotalQuizzesCompleted,
		&stats.TotalSimulatorGames,
		&stats.WinStreak,
		&stats.TotalWins,
		&stats.TotalLosses,
		&stats.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user stats not found")
	}

	return stats, err
}
