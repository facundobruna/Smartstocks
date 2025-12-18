package models

import (
	"database/sql"
	"time"
)

type Quiz struct {
	ID               string       `json:"id"`
	Difficulty       string       `json:"difficulty"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	PointsReward     int          `json:"points_reward"`
	TotalQuestions   int          `json:"total_questions"`
	TimeLimitMinutes int          `json:"time_limit_minutes"`
	IsActive         bool         `json:"is_active"`
	CreatedAt        time.Time    `json:"created_at"`
	ExpiresAt        sql.NullTime `json:"expires_at,omitempty"`
}

type QuizQuestion struct {
	ID            string    `json:"id"`
	QuizID        string    `json:"quiz_id"`
	QuestionText  string    `json:"question_text"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectOption string    `json:"-"` // No exponer en API
	Explanation   string    `json:"-"` // Solo mostrar después de responder
	Difficulty    string    `json:"difficulty"`
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"created_at"`
}

type QuizAttempt struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	QuizID           string    `json:"quiz_id"`
	Difficulty       string    `json:"difficulty"`
	Score            int       `json:"score"`
	TotalQuestions   int       `json:"total_questions"`
	CorrectAnswers   int       `json:"correct_answers"`
	PointsEarned     int       `json:"points_earned"`
	TimeTakenSeconds int       `json:"time_taken_seconds"`
	Answers          string    `json:"answers"` // JSON string
	StartedAt        time.Time `json:"started_at"`
	CompletedAt      time.Time `json:"completed_at"`
}

// DTOs

type GetQuizRequest struct {
	Difficulty string `uri:"difficulty" binding:"required,oneof=easy medium hard"`
}

type SubmitQuizRequest struct {
	QuizID           string                 `json:"quiz_id" binding:"required"`
	Answers          []QuizAnswerSubmission `json:"answers" binding:"required"`
	TimeTakenSeconds int                    `json:"time_taken_seconds"`
}

type QuizAnswerSubmission struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required,oneof=A B C D"`
}

type QuizResponse struct {
	Quiz       *Quiz          `json:"quiz"`
	Questions  []QuizQuestion `json:"questions"`
	CanAttempt bool           `json:"can_attempt"`
	Cooldown   *CooldownInfo  `json:"cooldown,omitempty"`
}

type CooldownInfo struct {
	NextAvailableAt time.Time `json:"next_available_at"`
	HoursRemaining  int       `json:"hours_remaining"`
}

type SubmitQuizResponse struct {
	AttemptID      string           `json:"attempt_id"`
	Score          int              `json:"score"`
	TotalQuestions int              `json:"total_questions"`
	CorrectAnswers int              `json:"correct_answers"`
	PointsEarned   int              `json:"points_earned"`
	NewTotalPoints int              `json:"new_total_points"`
	NewRank        string           `json:"new_rank"`
	Results        []QuestionResult `json:"results"`
}

type QuestionResult struct {
	QuestionID    string `json:"question_id"`
	QuestionText  string `json:"question_text"`
	YourAnswer    string `json:"your_answer"`
	CorrectAnswer string `json:"correct_answer"`
	IsCorrect     bool   `json:"is_correct"`
	Explanation   string `json:"explanation"`
}

type QuizHistoryResponse struct {
	Attempts []QuizAttempt `json:"attempts"`
	Stats    QuizStats     `json:"stats"`
}

type QuizStats struct {
	TotalAttempts   int     `json:"total_attempts"`
	EasyCompleted   int     `json:"easy_completed"`
	MediumCompleted int     `json:"medium_completed"`
	HardCompleted   int     `json:"hard_completed"`
	AverageScore    float64 `json:"average_score"`
	TotalPoints     int     `json:"total_points"`
}
