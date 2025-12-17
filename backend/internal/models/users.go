package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                string         `json:"id"`
	Username          string         `json:"username"`
	Email             string         `json:"email"`
	PasswordHash      string         `json:"-"`
	ProfilePictureURL sql.NullString `json:"profile_picture_url,omitempty"`
	SchoolID          sql.NullString `json:"school_id,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	LastLogin         sql.NullTime   `json:"last_login,omitempty"`
	EmailVerified     bool           `json:"email_verified"`
	VerificationToken sql.NullString `json:"-"`
	ResetToken        sql.NullString `json:"-"`
	ResetTokenExpires sql.NullTime   `json:"-"`
}

type UserStats struct {
	UserID                string    `json:"user_id"`
	Smartpoints           int       `json:"smartpoints"`
	RankTier              string    `json:"rank_tier"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	TotalSimulatorGames   int       `json:"total_simulator_games"`
	WinStreak             int       `json:"win_streak"`
	TotalWins             int       `json:"total_wins"`
	TotalLosses           int       `json:"total_losses"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type School struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username          string  `json:"username" binding:"required,min=3,max=50"`
	Email             string  `json:"email" binding:"required,email"`
	Password          string  `json:"password" binding:"required,min=8"`
	SchoolID          *string `json:"school_id,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         *UserInfo  `json:"user"`
	Stats        *UserStats `json:"stats"`
}

type UserInfo struct {
	ID                string  `json:"id"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	ProfilePictureURL *string `json:"profile_picture_url"`
	SchoolID          *string `json:"school_id"`
	EmailVerified     bool    `json:"email_verified"`
	CreatedAt         string  `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type EmailVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username          *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	SchoolID          *string `json:"school_id,omitempty"`
}
