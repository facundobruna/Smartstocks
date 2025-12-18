package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) CreateQuiz(quiz *models.Quiz) error {
	quiz.ID = uuid.New().String()
	quiz.CreatedAt = time.Now()

	query := `
		INSERT INTO quizzes (id, difficulty, title, description, points_reward, 
							total_questions, time_limit_minutes, is_active, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		quiz.ID,
		quiz.Difficulty,
		quiz.Title,
		quiz.Description,
		quiz.PointsReward,
		quiz.TotalQuestions,
		quiz.TimeLimitMinutes,
		quiz.IsActive,
		quiz.ExpiresAt,
	)

	return err
}

func (r *QuizRepository) CreateQuestion(question *models.QuizQuestion) error {
	question.ID = uuid.New().String()
	question.CreatedAt = time.Now()

	query := `
		INSERT INTO quiz_questions (id, quiz_id, question_text, option_a, option_b, 
									option_c, option_d, correct_option, explanation, 
									difficulty, category)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		question.ID,
		question.QuizID,
		question.QuestionText,
		question.OptionA,
		question.OptionB,
		question.OptionC,
		question.OptionD,
		question.CorrectOption,
		question.Explanation,
		question.Difficulty,
		question.Category,
	)

	return err
}

func (r *QuizRepository) GetActiveQuizByDifficulty(difficulty string) (*models.Quiz, error) {
	quiz := &models.Quiz{}
	query := `
		SELECT id, difficulty, title, description, points_reward, total_questions,
			   time_limit_minutes, is_active, created_at, expires_at
		FROM quizzes
		WHERE difficulty = ? AND is_active = TRUE
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, difficulty).Scan(
		&quiz.ID,
		&quiz.Difficulty,
		&quiz.Title,
		&quiz.Description,
		&quiz.PointsReward,
		&quiz.TotalQuestions,
		&quiz.TimeLimitMinutes,
		&quiz.IsActive,
		&quiz.CreatedAt,
		&quiz.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return quiz, err
}

func (r *QuizRepository) GetQuestionsByQuizID(quizID string) ([]models.QuizQuestion, error) {
	query := `
		SELECT id, quiz_id, question_text, option_a, option_b, option_c, option_d,
			   correct_option, explanation, difficulty, category, created_at
		FROM quiz_questions
		WHERE quiz_id = ?
		ORDER BY created_at
	`

	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.QuizQuestion
	for rows.Next() {
		var q models.QuizQuestion
		err := rows.Scan(
			&q.ID,
			&q.QuizID,
			&q.QuestionText,
			&q.OptionA,
			&q.OptionB,
			&q.OptionC,
			&q.OptionD,
			&q.CorrectOption,
			&q.Explanation,
			&q.Difficulty,
			&q.Category,
			&q.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}

func (r *QuizRepository) CreateAttempt(attempt *models.QuizAttempt) error {
	attempt.ID = uuid.New().String()
	attempt.StartedAt = time.Now()
	attempt.CompletedAt = time.Now()

	query := `
		INSERT INTO quiz_attempts (id, user_id, quiz_id, difficulty, score, 
								  total_questions, correct_answers, points_earned,
								  time_taken_seconds, answers)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		attempt.ID,
		attempt.UserID,
		attempt.QuizID,
		attempt.Difficulty,
		attempt.Score,
		attempt.TotalQuestions,
		attempt.CorrectAnswers,
		attempt.PointsEarned,
		attempt.TimeTakenSeconds,
		attempt.Answers,
	)

	return err
}

func (r *QuizRepository) CheckCooldown(userID, difficulty string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM daily_quiz_cooldowns
		WHERE user_id = ? AND difficulty = ? AND last_attempt_date = CURDATE()
	`

	err := r.db.QueryRow(query, userID, difficulty).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *QuizRepository) SetCooldown(userID, difficulty string) error {
	query := `
		INSERT INTO daily_quiz_cooldowns (id, user_id, difficulty, last_attempt_date)
		VALUES (?, ?, ?, CURDATE())
		ON DUPLICATE KEY UPDATE last_attempt_date = CURDATE()
	`

	_, err := r.db.Exec(query, uuid.New().String(), userID, difficulty)
	return err
}

func (r *QuizRepository) UpdateUserStatsAfterQuiz(userID string, pointsEarned int) error {
	query := `CALL update_user_stats_after_quiz(?, ?)`
	_, err := r.db.Exec(query, userID, pointsEarned)
	return err
}

func (r *QuizRepository) GetUserAttempts(userID string, limit int) ([]models.QuizAttempt, error) {
	query := `
		SELECT id, user_id, quiz_id, difficulty, score, total_questions,
			   correct_answers, points_earned, time_taken_seconds, answers,
			   started_at, completed_at
		FROM quiz_attempts
		WHERE user_id = ?
		ORDER BY completed_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.QuizAttempt
	for rows.Next() {
		var a models.QuizAttempt
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.QuizID,
			&a.Difficulty,
			&a.Score,
			&a.TotalQuestions,
			&a.CorrectAnswers,
			&a.PointsEarned,
			&a.TimeTakenSeconds,
			&a.Answers,
			&a.StartedAt,
			&a.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}

	return attempts, nil
}

func (r *QuizRepository) GetQuizStats(userID string) (*models.QuizStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_attempts,
			SUM(CASE WHEN difficulty = 'easy' THEN 1 ELSE 0 END) as easy_completed,
			SUM(CASE WHEN difficulty = 'medium' THEN 1 ELSE 0 END) as medium_completed,
			SUM(CASE WHEN difficulty = 'hard' THEN 1 ELSE 0 END) as hard_completed,
			AVG(score) as average_score,
			SUM(points_earned) as total_points
		FROM quiz_attempts
		WHERE user_id = ?
	`

	stats := &models.QuizStats{}
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalAttempts,
		&stats.EasyCompleted,
		&stats.MediumCompleted,
		&stats.HardCompleted,
		&stats.AverageScore,
		&stats.TotalPoints,
	)

	if err == sql.ErrNoRows {
		return &models.QuizStats{}, nil
	}

	return stats, err
}

func (r *QuizRepository) GetQuestionByID(questionID string) (*models.QuizQuestion, error) {
	question := &models.QuizQuestion{}
	query := `
		SELECT id, quiz_id, question_text, option_a, option_b, option_c, option_d,
			   correct_option, explanation, difficulty, category, created_at
		FROM quiz_questions
		WHERE id = ?
	`

	err := r.db.QueryRow(query, questionID).Scan(
		&question.ID,
		&question.QuizID,
		&question.QuestionText,
		&question.OptionA,
		&question.OptionB,
		&question.OptionC,
		&question.OptionD,
		&question.CorrectOption,
		&question.Explanation,
		&question.Difficulty,
		&question.Category,
		&question.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("question not found")
	}

	return question, err
}
