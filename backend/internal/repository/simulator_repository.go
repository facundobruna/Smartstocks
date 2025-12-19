package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type SimulatorRepository struct {
	db *sql.DB
}

func NewSimulatorRepository(db *sql.DB) *SimulatorRepository {
	return &SimulatorRepository{db: db}
}

// CreateScenario crea un nuevo escenario
func (r *SimulatorRepository) CreateScenario(scenario *models.SimulatorScenario) error {
	scenario.ID = uuid.New().String()

	chartDataJSON, err := json.Marshal(scenario.ChartData)
	if err != nil {
		return fmt.Errorf("error marshaling chart data: %w", err)
	}

	query := `
		INSERT INTO simulator_scenarios (
			id, difficulty, news_content, chart_data, 
			correct_decision, explanation, created_at, expires_at, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		scenario.ID,
		scenario.Difficulty,
		scenario.NewsContent,
		chartDataJSON,
		scenario.CorrectDecision,
		scenario.Explanation,
		scenario.CreatedAt,
		scenario.ExpiresAt,
		scenario.IsActive,
	)

	return err
}

// GetActiveScenarioByDifficulty obtiene un escenario activo por dificultad
func (r *SimulatorRepository) GetActiveScenarioByDifficulty(difficulty models.SimulatorDifficulty) (*models.SimulatorScenario, error) {
	scenario := &models.SimulatorScenario{}
	var chartDataJSON []byte

	query := `
		SELECT id, difficulty, news_content, chart_data, 
			   correct_decision, explanation, created_at, expires_at, is_active
		FROM simulator_scenarios
		WHERE difficulty = ? AND is_active = TRUE AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, difficulty).Scan(
		&scenario.ID,
		&scenario.Difficulty,
		&scenario.NewsContent,
		&chartDataJSON,
		&scenario.CorrectDecision,
		&scenario.Explanation,
		&scenario.CreatedAt,
		&scenario.ExpiresAt,
		&scenario.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No hay escenarios disponibles
	}
	if err != nil {
		return nil, err
	}

	// Parsear chart data JSON
	if err := json.Unmarshal(chartDataJSON, &scenario.ChartData); err != nil {
		return nil, fmt.Errorf("error unmarshaling chart data: %w", err)
	}

	return scenario, nil
}

// GetScenarioByID obtiene un escenario por ID
func (r *SimulatorRepository) GetScenarioByID(scenarioID string) (*models.SimulatorScenario, error) {
	scenario := &models.SimulatorScenario{}
	var chartDataJSON []byte

	query := `
		SELECT id, difficulty, news_content, chart_data, 
			   correct_decision, explanation, created_at, expires_at, is_active
		FROM simulator_scenarios
		WHERE id = ?
	`

	err := r.db.QueryRow(query, scenarioID).Scan(
		&scenario.ID,
		&scenario.Difficulty,
		&scenario.NewsContent,
		&chartDataJSON,
		&scenario.CorrectDecision,
		&scenario.Explanation,
		&scenario.CreatedAt,
		&scenario.ExpiresAt,
		&scenario.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("scenario not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(chartDataJSON, &scenario.ChartData); err != nil {
		return nil, fmt.Errorf("error unmarshaling chart data: %w", err)
	}

	return scenario, nil
}

// CheckCooldown verifica si el usuario puede intentar un quiz de cierta dificultad
func (r *SimulatorRepository) CheckCooldown(userID string, difficulty models.SimulatorDifficulty) (bool, error) {
	var canAttempt bool

	query := `CALL check_simulator_cooldown(?, ?, @can_attempt)`
	_, err := r.db.Exec(query, userID, difficulty)
	if err != nil {
		return false, err
	}

	err = r.db.QueryRow(`SELECT @can_attempt`).Scan(&canAttempt)
	if err != nil {
		return false, err
	}

	return canAttempt, nil
}

// GetLastCooldown obtiene el último cooldown del usuario para una dificultad
func (r *SimulatorRepository) GetLastCooldown(userID string, difficulty models.SimulatorDifficulty) (*models.DailySimulatorCooldown, error) {
	cooldown := &models.DailySimulatorCooldown{}

	query := `
		SELECT id, user_id, difficulty, last_attempt_date, attempts_count, created_at, updated_at
		FROM daily_simulator_cooldowns
		WHERE user_id = ? AND difficulty = ?
		ORDER BY last_attempt_date DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, userID, difficulty).Scan(
		&cooldown.ID,
		&cooldown.UserID,
		&cooldown.Difficulty,
		&cooldown.LastAttemptDate,
		&cooldown.AttemptsCount,
		&cooldown.CreatedAt,
		&cooldown.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return cooldown, nil
}

// RecordAttempt registra un intento de simulador
func (r *SimulatorRepository) RecordAttempt(attempt *models.SimulatorAttempt) error {
	timeTaken := sql.NullInt64{Valid: false}
	if attempt.TimeTakenSeconds.Valid {
		timeTaken = attempt.TimeTakenSeconds
	}

	query := `
		CALL record_simulator_attempt(?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		attempt.UserID,
		attempt.ScenarioID,
		attempt.Difficulty,
		attempt.UserDecision,
		attempt.WasCorrect,
		attempt.PointsEarned,
		timeTaken,
	)

	return err
}

// GetUserAttempts obtiene el historial de intentos del usuario
func (r *SimulatorRepository) GetUserAttempts(userID string, limit int) ([]models.SimulatorAttemptWithDetails, error) {
	query := `
		SELECT 
			a.id, a.user_id, a.scenario_id, a.difficulty,
			a.user_decision, a.was_correct, a.points_earned,
			a.time_taken_seconds, a.created_at,
			s.news_content, s.explanation
		FROM simulator_attempts a
		JOIN simulator_scenarios s ON a.scenario_id = s.id
		WHERE a.user_id = ?
		ORDER BY a.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.SimulatorAttemptWithDetails
	for rows.Next() {
		var attempt models.SimulatorAttemptWithDetails
		err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.ScenarioID,
			&attempt.Difficulty,
			&attempt.UserDecision,
			&attempt.WasCorrect,
			&attempt.PointsEarned,
			&attempt.TimeTakenSeconds,
			&attempt.CreatedAt,
			&attempt.NewsContent,
			&attempt.Explanation,
		)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

// GetUserStats obtiene estadísticas del simulador para un usuario
func (r *SimulatorRepository) GetUserStats(userID string) (*models.SimulatorStats, error) {
	// Estadísticas generales
	var totalAttempts, correctAttempts, totalPoints int
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN was_correct THEN 1 ELSE 0 END) as correct,
			SUM(CASE WHEN was_correct THEN points_earned ELSE 0 END) as points
		FROM simulator_attempts
		WHERE user_id = ?
	`

	err := r.db.QueryRow(query, userID).Scan(&totalAttempts, &correctAttempts, &totalPoints)
	if err != nil {
		return nil, err
	}

	// Estadísticas por dificultad
	byDifficulty := make(map[string]models.SimulatorDifficultyStats)

	difficulties := []models.SimulatorDifficulty{
		models.SimulatorDifficultyEasy,
		models.SimulatorDifficultyMedium,
		models.SimulatorDifficultyHard,
	}

	for _, diff := range difficulties {
		var attempts, correct, points int
		query := `
			SELECT 
				COUNT(*) as total,
				SUM(CASE WHEN was_correct THEN 1 ELSE 0 END) as correct,
				SUM(CASE WHEN was_correct THEN points_earned ELSE 0 END) as points
			FROM simulator_attempts
			WHERE user_id = ? AND difficulty = ?
		`

		err := r.db.QueryRow(query, userID, diff).Scan(&attempts, &correct, &points)
		if err != nil {
			return nil, err
		}

		accuracy := 0.0
		if attempts > 0 {
			accuracy = float64(correct) / float64(attempts) * 100
		}

		byDifficulty[string(diff)] = models.SimulatorDifficultyStats{
			Attempts:     attempts,
			Correct:      correct,
			AccuracyRate: accuracy,
			PointsEarned: points,
		}
	}

	// Calcular accuracy general
	accuracyRate := 0.0
	if totalAttempts > 0 {
		accuracyRate = float64(correctAttempts) / float64(totalAttempts) * 100
	}

	stats := &models.SimulatorStats{
		TotalAttempts:   totalAttempts,
		CorrectAttempts: correctAttempts,
		AccuracyRate:    accuracyRate,
		TotalPoints:     totalPoints,
		ByDifficulty:    byDifficulty,
	}

	return stats, nil
}

// CleanupExpiredScenarios limpia escenarios expirados
func (r *SimulatorRepository) CleanupExpiredScenarios() error {
	query := `CALL cleanup_expired_scenarios()`
	_, err := r.db.Exec(query)
	return err
}

// DeactivateScenario desactiva un escenario
func (r *SimulatorRepository) DeactivateScenario(scenarioID string) error {
	query := `UPDATE simulator_scenarios SET is_active = FALSE WHERE id = ?`
	_, err := r.db.Exec(query, scenarioID)
	return err
}

// GetScenarioUsageCount obtiene cuántas veces se ha usado un escenario
func (r *SimulatorRepository) GetScenarioUsageCount(scenarioID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM simulator_attempts WHERE scenario_id = ?`
	err := r.db.QueryRow(query, scenarioID).Scan(&count)
	return count, err
}
