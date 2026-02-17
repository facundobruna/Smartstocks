package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type PvPRepository struct {
	db *sql.DB
}

func NewPvPRepository(db *sql.DB) *PvPRepository {
	return &PvPRepository{db: db}
}

// === QUEUE MANAGEMENT ===

// JoinQueue añade un usuario a la cola
func (r *PvPRepository) JoinQueue(userID, rankTier string) (*models.PvPQueueEntry, error) {
	entry := &models.PvPQueueEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		RankTier:  rankTier,
		JoinedAt:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute), // Expira en 5 minutos
		IsActive:  true,
	}

	query := `
		INSERT INTO pvp_queue (id, user_id, rank_tier, joined_at, expires_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			joined_at = VALUES(joined_at),
			expires_at = VALUES(expires_at),
			is_active = TRUE
	`

	_, err := r.db.Exec(query,
		entry.ID,
		entry.UserID,
		entry.RankTier,
		entry.JoinedAt,
		entry.ExpiresAt,
		entry.IsActive,
	)

	return entry, err
}

// LeaveQueue saca un usuario de la cola
func (r *PvPRepository) LeaveQueue(userID string) error {
	query := `UPDATE pvp_queue SET is_active = FALSE WHERE user_id = ? AND is_active = TRUE`
	_, err := r.db.Exec(query, userID)
	return err
}

// FindOpponent busca un oponente en la cola
func (r *PvPRepository) FindOpponent(userID, rankTier string) (*models.PvPQueueEntry, error) {
	var opponentID string

	// Llamar al stored procedure
	query := `CALL find_opponent(?, ?, @opponent_id)`
	_, err := r.db.Exec(query, userID, rankTier)
	if err != nil {
		return nil, err
	}

	// Obtener el resultado
	err = r.db.QueryRow(`SELECT @opponent_id`).Scan(&opponentID)
	if err != nil {
		return nil, err
	}

	if opponentID == "" {
		return nil, nil // No hay oponente disponible
	}

	// Obtener los detalles del oponente
	opponent := &models.PvPQueueEntry{}
	query = `
		SELECT id, user_id, rank_tier, joined_at, expires_at, is_active
		FROM pvp_queue
		WHERE user_id = ? AND is_active = TRUE
	`

	err = r.db.QueryRow(query, opponentID).Scan(
		&opponent.ID,
		&opponent.UserID,
		&opponent.RankTier,
		&opponent.JoinedAt,
		&opponent.ExpiresAt,
		&opponent.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return opponent, err
}

// GetQueuePosition obtiene la posición en la cola
func (r *PvPRepository) GetQueuePosition(userID string) (int, error) {
	var position int
	query := `
		SELECT COUNT(*) + 1
		FROM pvp_queue
		WHERE is_active = TRUE
		AND joined_at < (SELECT joined_at FROM pvp_queue WHERE user_id = ? AND is_active = TRUE)
	`

	err := r.db.QueryRow(query, userID).Scan(&position)
	return position, err
}

// CleanupExpiredQueue limpia entradas expiradas de la cola
func (r *PvPRepository) CleanupExpiredQueue() error {
	query := `CALL cleanup_expired_queue()`
	_, err := r.db.Exec(query)
	return err
}

// === MATCH MANAGEMENT ===

// CreateMatch crea una nueva partida
func (r *PvPRepository) CreateMatch(player1ID, player2ID string) (*models.PvPMatch, error) {
	match := &models.PvPMatch{
		ID:           uuid.New().String(),
		Player1ID:    player1ID,
		Player2ID:    player2ID,
		Player1Score: 0,
		Player2Score: 0,
		Status:       models.PvPMatchStatusWaiting,
		CurrentRound: 0,
		TotalRounds:  5,
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO pvp_matches (
			id, player1_id, player2_id, player1_score, player2_score,
			status, current_round, total_rounds, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		match.ID,
		match.Player1ID,
		match.Player2ID,
		match.Player1Score,
		match.Player2Score,
		match.Status,
		match.CurrentRound,
		match.TotalRounds,
		match.CreatedAt,
	)

	return match, err
}

// GetMatchByID obtiene una partida por ID
func (r *PvPRepository) GetMatchByID(matchID string) (*models.PvPMatch, error) {
	match := &models.PvPMatch{}

	query := `
		SELECT id, player1_id, player2_id, player1_score, player2_score,
			   winner_id, status, current_round, total_rounds,
			   started_at, completed_at, created_at
		FROM pvp_matches
		WHERE id = ?
	`

	err := r.db.QueryRow(query, matchID).Scan(
		&match.ID,
		&match.Player1ID,
		&match.Player2ID,
		&match.Player1Score,
		&match.Player2Score,
		&match.WinnerID,
		&match.Status,
		&match.CurrentRound,
		&match.TotalRounds,
		&match.StartedAt,
		&match.CompletedAt,
		&match.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("match not found")
	}

	return match, err
}

// StartMatch marca una partida como iniciada
func (r *PvPRepository) StartMatch(matchID string) error {
	query := `
		UPDATE pvp_matches 
		SET status = ?, started_at = ?, current_round = 1
		WHERE id = ?
	`
	_, err := r.db.Exec(query, models.PvPMatchStatusInProgress, time.Now(), matchID)
	return err
}

// UpdateMatchScores actualiza los puntajes de una partida
func (r *PvPRepository) UpdateMatchScores(matchID string, player1Score, player2Score int) error {
	query := `
		UPDATE pvp_matches 
		SET player1_score = ?, player2_score = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, player1Score, player2Score, matchID)
	return err
}

// CompleteMatch marca una partida como completada
func (r *PvPRepository) CompleteMatch(matchID, winnerID string) error {
	query := `
		UPDATE pvp_matches 
		SET status = ?, winner_id = ?, completed_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, models.PvPMatchStatusCompleted, winnerID, time.Now(), matchID)
	return err
}

// === ROUND MANAGEMENT ===

// CreateRound crea una nueva ronda
func (r *PvPRepository) CreateRound(matchID string, roundNumber int, scenarioID string, correctDecision models.SimulatorDecision) (*models.PvPRound, error) {
	round := &models.PvPRound{
		ID:              uuid.New().String(),
		MatchID:         matchID,
		RoundNumber:     roundNumber,
		ScenarioID:      scenarioID,
		CorrectDecision: correctDecision,
		StartedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		CreatedAt:       time.Now(),
	}

	query := `
		INSERT INTO pvp_rounds (
			id, match_id, round_number, scenario_id, correct_decision,
			started_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		round.ID,
		round.MatchID,
		round.RoundNumber,
		round.ScenarioID,
		round.CorrectDecision,
		round.StartedAt,
		round.CreatedAt,
	)

	return round, err
}

// GetRound obtiene una ronda específica
func (r *PvPRepository) GetRound(matchID string, roundNumber int) (*models.PvPRound, error) {
	round := &models.PvPRound{}

	query := `
		SELECT id, match_id, round_number, scenario_id,
			   player1_decision, player2_decision,
			   player1_time_seconds, player2_time_seconds,
			   player1_correct, player2_correct,
			   player1_points, player2_points,
			   correct_decision, started_at, completed_at, created_at
		FROM pvp_rounds
		WHERE match_id = ? AND round_number = ?
	`

	err := r.db.QueryRow(query, matchID, roundNumber).Scan(
		&round.ID,
		&round.MatchID,
		&round.RoundNumber,
		&round.ScenarioID,
		&round.Player1Decision,
		&round.Player2Decision,
		&round.Player1TimeSeconds,
		&round.Player2TimeSeconds,
		&round.Player1Correct,
		&round.Player2Correct,
		&round.Player1Points,
		&round.Player2Points,
		&round.CorrectDecision,
		&round.StartedAt,
		&round.CompletedAt,
		&round.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("round not found")
	}

	return round, err
}

// SubmitRoundDecision registra la decisión de un jugador
func (r *PvPRepository) SubmitRoundDecision(matchID string, roundNumber int, playerID string, decision models.SimulatorDecision, timeElapsed float64) error {
	// Determinar si es player1 o player2
	match, err := r.GetMatchByID(matchID)
	if err != nil {
		return err
	}

	var query string
	if playerID == match.Player1ID {
		query = `
			UPDATE pvp_rounds
			SET player1_decision = ?, player1_time_seconds = ?
			WHERE match_id = ? AND round_number = ?
		`
	} else {
		query = `
			UPDATE pvp_rounds
			SET player2_decision = ?, player2_time_seconds = ?
			WHERE match_id = ? AND round_number = ?
		`
	}

	_, err = r.db.Exec(query, decision, timeElapsed, matchID, roundNumber)
	return err
}

// CompleteRound marca una ronda como completada y calcula puntos
func (r *PvPRepository) CompleteRound(matchID string, roundNumber int) error {
	round, err := r.GetRound(matchID, roundNumber)
	if err != nil {
		return err
	}

	// Verificar que ambos jugadores hayan decidido
	if !round.Player1Decision.Valid || !round.Player2Decision.Valid {
		return errors.New("both players must submit decisions")
	}

	// Calcular si cada jugador acertó
	player1Correct := models.SimulatorDecision(round.Player1Decision.String) == round.CorrectDecision
	player2Correct := models.SimulatorDecision(round.Player2Decision.String) == round.CorrectDecision

	// Calcular puntos
	player1Points := 0
	player2Points := 0

	if player1Correct {
		player1Points = models.CalculateRoundPoints(true, round.Player1TimeSeconds.Float64)
	}
	if player2Correct {
		player2Points = models.CalculateRoundPoints(true, round.Player2TimeSeconds.Float64)
	}

	// Actualizar ronda
	query := `
		UPDATE pvp_rounds
		SET player1_correct = ?, player2_correct = ?,
			player1_points = ?, player2_points = ?,
			completed_at = ?
		WHERE match_id = ? AND round_number = ?
	`

	_, err = r.db.Exec(query,
		player1Correct, player2Correct,
		player1Points, player2Points,
		time.Now(),
		matchID, roundNumber,
	)

	return err
}

// GetMatchRounds obtiene todas las rondas de una partida
func (r *PvPRepository) GetMatchRounds(matchID string) ([]models.PvPRound, error) {
	query := `
		SELECT id, match_id, round_number, scenario_id,
			   player1_decision, player2_decision,
			   player1_time_seconds, player2_time_seconds,
			   player1_correct, player2_correct,
			   player1_points, player2_points,
			   correct_decision, started_at, completed_at, created_at
		FROM pvp_rounds
		WHERE match_id = ?
		ORDER BY round_number ASC
	`

	rows, err := r.db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rounds []models.PvPRound
	for rows.Next() {
		var round models.PvPRound
		err := rows.Scan(
			&round.ID,
			&round.MatchID,
			&round.RoundNumber,
			&round.ScenarioID,
			&round.Player1Decision,
			&round.Player2Decision,
			&round.Player1TimeSeconds,
			&round.Player2TimeSeconds,
			&round.Player1Correct,
			&round.Player2Correct,
			&round.Player1Points,
			&round.Player2Points,
			&round.CorrectDecision,
			&round.StartedAt,
			&round.CompletedAt,
			&round.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rounds = append(rounds, round)
	}

	return rounds, nil
}

// === STATS & HISTORY ===

// GetUserMatches obtiene el historial de partidas de un usuario
func (r *PvPRepository) GetUserMatches(userID string, limit int) ([]models.PvPMatch, error) {
	query := `
		SELECT id, player1_id, player2_id, player1_score, player2_score,
			   winner_id, status, current_round, total_rounds,
			   started_at, completed_at, created_at
		FROM pvp_matches
		WHERE (player1_id = ? OR player2_id = ?)
		AND status = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, userID, userID, models.PvPMatchStatusCompleted, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.PvPMatch
	for rows.Next() {
		var match models.PvPMatch
		err := rows.Scan(
			&match.ID,
			&match.Player1ID,
			&match.Player2ID,
			&match.Player1Score,
			&match.Player2Score,
			&match.WinnerID,
			&match.Status,
			&match.CurrentRound,
			&match.TotalRounds,
			&match.StartedAt,
			&match.CompletedAt,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

// GetUserPvPStats obtiene estadísticas PvP del usuario
func (r *PvPRepository) GetUserPvPStats(userID string) (*models.PvPStats, error) {
	stats := &models.PvPStats{}

	// Estadísticas básicas
	query := `
		SELECT 
			COALESCE(COUNT(*), 0) as total,
			COALESCE(SUM(CASE WHEN winner_id = ? THEN 1 ELSE 0 END), 0) as wins,
			COALESCE(SUM(CASE WHEN winner_id IS NOT NULL AND winner_id != ? THEN 1 ELSE 0 END), 0) as losses,
			COALESCE(SUM(CASE 
				WHEN winner_id IS NULL 
				AND status = ? 
				AND player1_score = player2_score 
				THEN 1 ELSE 0 
			END), 0) as ties
		FROM pvp_matches
		WHERE (player1_id = ? OR player2_id = ?)
		AND status = ?
	`

	err := r.db.QueryRow(query, userID, userID, models.PvPMatchStatusCompleted, userID, userID, models.PvPMatchStatusCompleted).Scan(
		&stats.TotalMatches,
		&stats.Wins,
		&stats.Losses,
		&stats.Ties,
	)
	if err != nil {
		return nil, err
	}

	// Calcular win rate
	if stats.TotalMatches > 0 {
		stats.WinRate = float64(stats.Wins) / float64(stats.TotalMatches) * 100
	}

	// Obtener racha actual de user_stats
	query = `SELECT win_streak FROM user_stats WHERE user_id = ?`
	err = r.db.QueryRow(query, userID).Scan(&stats.CurrentStreak)
	if err != nil {
		return nil, err
	}

	// Mejor racha histórica (aproximada por total_wins)
	query = `SELECT total_wins FROM user_stats WHERE user_id = ?`
	err = r.db.QueryRow(query, userID).Scan(&stats.BestStreak)
	if err != nil {
		return nil, err
	}

	// Puntos ganados vs perdidos
	query = `
		SELECT 
			COALESCE(SUM(CASE 
				WHEN winner_id = ? THEN 
					CASE WHEN player1_id = ? THEN player1_score ELSE player2_score END
				ELSE 0 
			END), 0) as won,
			COALESCE(SUM(CASE 
				WHEN winner_id != ? AND winner_id IS NOT NULL THEN 100
				ELSE 0 
			END), 0) as lost
		FROM pvp_matches
		WHERE (player1_id = ? OR player2_id = ?)
		AND status = ?
	`

	err = r.db.QueryRow(query, userID, userID, userID, userID, userID, models.PvPMatchStatusCompleted).Scan(
		&stats.TotalPointsWon,
		&stats.TotalPointsLost,
	)

	return stats, err
}

// UpdatePvPStats actualiza las estadísticas después de una partida
func (r *PvPRepository) UpdatePvPStats(winnerID, loserID string, winnerPoints int) error {
	query := `CALL update_pvp_stats(?, ?, ?, ?)`
	_, err := r.db.Exec(query, winnerID, loserID, winnerPoints, true)
	return err
}
