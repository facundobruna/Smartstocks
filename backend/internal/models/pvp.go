package models

import (
	"database/sql"
	"time"
)

// PvPMatchStatus representa los estados de una partida
type PvPMatchStatus string

const (
	PvPMatchStatusWaiting    PvPMatchStatus = "waiting"
	PvPMatchStatusInProgress PvPMatchStatus = "in_progress"
	PvPMatchStatusCompleted  PvPMatchStatus = "completed"
	PvPMatchStatusCancelled  PvPMatchStatus = "cancelled"
)

// PvPMatch representa una partida PvP
type PvPMatch struct {
	ID           string         `json:"id"`
	Player1ID    string         `json:"player1_id"`
	Player2ID    string         `json:"player2_id"`
	Player1Score int            `json:"player1_score"`
	Player2Score int            `json:"player2_score"`
	WinnerID     sql.NullString `json:"winner_id,omitempty"`
	Status       PvPMatchStatus `json:"status"`
	CurrentRound int            `json:"current_round"`
	TotalRounds  int            `json:"total_rounds"`
	StartedAt    sql.NullTime   `json:"started_at,omitempty"`
	CompletedAt  sql.NullTime   `json:"completed_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// PvPRound representa una ronda de una partida
type PvPRound struct {
	ID                 string            `json:"id"`
	MatchID            string            `json:"match_id"`
	RoundNumber        int               `json:"round_number"`
	ScenarioID         string            `json:"scenario_id"`
	Player1Decision    sql.NullString    `json:"player1_decision,omitempty"`
	Player2Decision    sql.NullString    `json:"player2_decision,omitempty"`
	Player1TimeSeconds sql.NullFloat64   `json:"player1_time_seconds,omitempty"`
	Player2TimeSeconds sql.NullFloat64   `json:"player2_time_seconds,omitempty"`
	Player1Correct     sql.NullBool      `json:"player1_correct,omitempty"`
	Player2Correct     sql.NullBool      `json:"player2_correct,omitempty"`
	Player1Points      int               `json:"player1_points"`
	Player2Points      int               `json:"player2_points"`
	CorrectDecision    SimulatorDecision `json:"-"` // No enviar hasta que termine
	StartedAt          sql.NullTime      `json:"started_at,omitempty"`
	CompletedAt        sql.NullTime      `json:"completed_at,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

// PvPQueueEntry representa una entrada en la cola de matchmaking
type PvPQueueEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RankTier  string    `json:"rank_tier"`
	JoinedAt  time.Time `json:"joined_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

// === REQUEST/RESPONSE MODELS ===

// JoinQueueRequest representa la solicitud para unirse a la cola
type JoinQueueRequest struct {
	// No necesita parámetros, usa el user_id del token
}

// JoinQueueResponse representa la respuesta al unirse a la cola
type JoinQueueResponse struct {
	QueueID   string    `json:"queue_id"`
	Position  int       `json:"position"`
	ExpiresAt time.Time `json:"expires_at"`
	Message   string    `json:"message"`
}

// LeaveQueueRequest representa la solicitud para salir de la cola
type LeaveQueueRequest struct {
	// No necesita parámetros
}

// MatchFoundResponse representa cuando se encuentra un oponente
type MatchFoundResponse struct {
	MatchID     string    `json:"match_id"`
	OpponentID  string    `json:"opponent_id"`
	Opponent    *UserInfo `json:"opponent"`
	TotalRounds int       `json:"total_rounds"`
	Message     string    `json:"message"`
}

// RoundStartResponse representa el inicio de una ronda
type RoundStartResponse struct {
	MatchID     string                    `json:"match_id"`
	RoundNumber int                       `json:"round_number"`
	TotalRounds int                       `json:"total_rounds"`
	Scenario    SimulatorScenarioResponse `json:"scenario"`
	TimeLimit   int                       `json:"time_limit_seconds"` // 15 segundos
}

// SubmitPvPDecisionRequest representa el envío de una decisión en PvP
type SubmitPvPDecisionRequest struct {
	MatchID     string            `json:"match_id" binding:"required"`
	RoundNumber int               `json:"round_number" binding:"required,min=1"`
	Decision    SimulatorDecision `json:"decision" binding:"required,oneof=buy sell hold"`
	TimeElapsed float64           `json:"time_elapsed" binding:"required,min=0,max=15"`
}

// RoundResultResponse representa el resultado de una ronda
type RoundResultResponse struct {
	MatchID            string            `json:"match_id"`
	RoundNumber        int               `json:"round_number"`
	YourDecision       SimulatorDecision `json:"your_decision"`
	OpponentDecision   SimulatorDecision `json:"opponent_decision"`
	CorrectDecision    SimulatorDecision `json:"correct_decision"`
	YourCorrect        bool              `json:"your_correct"`
	OpponentCorrect    bool              `json:"opponent_correct"`
	YourTime           float64           `json:"your_time"`
	OpponentTime       float64           `json:"opponent_time"`
	YourPoints         int               `json:"your_points"`
	OpponentPoints     int               `json:"opponent_points"`
	YourTotalScore     int               `json:"your_total_score"`
	OpponentTotalScore int               `json:"opponent_total_score"`
	Explanation        string            `json:"explanation"`
	IsMatchComplete    bool              `json:"is_match_complete"`
}

// MatchResultResponse representa el resultado final de la partida
type MatchResultResponse struct {
	MatchID            string         `json:"match_id"`
	Winner             string         `json:"winner"` // "you", "opponent", "tie"
	YourFinalScore     int            `json:"your_final_score"`
	OpponentFinalScore int            `json:"opponent_final_score"`
	PointsGained       int            `json:"points_gained"` // Puede ser negativo
	NewTotalPoints     int            `json:"new_total_points"`
	NewRankTier        string         `json:"new_rank_tier"`
	WinStreak          int            `json:"win_streak"`
	StreakBonus        int            `json:"streak_bonus,omitempty"`
	Rounds             []RoundSummary `json:"rounds"`
}

// RoundSummary resumen de una ronda para el resultado final
type RoundSummary struct {
	RoundNumber      int               `json:"round_number"`
	YourDecision     SimulatorDecision `json:"your_decision"`
	OpponentDecision SimulatorDecision `json:"opponent_decision"`
	CorrectDecision  SimulatorDecision `json:"correct_decision"`
	YourPoints       int               `json:"your_points"`
	OpponentPoints   int               `json:"opponent_points"`
}

// PvPHistoryResponse representa el historial de partidas
type PvPHistoryResponse struct {
	Matches []PvPMatchWithDetails `json:"matches"`
	Stats   PvPStats              `json:"stats"`
}

// PvPMatchWithDetails incluye detalles de jugadores
type PvPMatchWithDetails struct {
	PvPMatch
	Player1 *UserInfo  `json:"player1"`
	Player2 *UserInfo  `json:"player2"`
	Winner  *UserInfo  `json:"winner,omitempty"`
	Rounds  []PvPRound `json:"rounds,omitempty"`
}

// PvPStats representa estadísticas de PvP
type PvPStats struct {
	TotalMatches    int     `json:"total_matches"`
	Wins            int     `json:"wins"`
	Losses          int     `json:"losses"`
	Ties            int     `json:"ties"`
	WinRate         float64 `json:"win_rate"`
	CurrentStreak   int     `json:"current_streak"`
	BestStreak      int     `json:"best_streak"`
	TotalPointsWon  int     `json:"total_points_won"`
	TotalPointsLost int     `json:"total_points_lost"`
}

// === WEBSOCKET MESSAGES ===

// WSMessageType tipos de mensajes WebSocket
type WSMessageType string

const (
	WSMsgTypeMatchFound   WSMessageType = "match_found"
	WSMsgTypeRoundStart   WSMessageType = "round_start"
	WSMsgTypeRoundResult  WSMessageType = "round_result"
	WSMsgTypeMatchResult  WSMessageType = "match_result"
	WSMsgTypeError        WSMessageType = "error"
	WSMsgTypeOpponentLeft WSMessageType = "opponent_left"
	WSMsgTypePing         WSMessageType = "ping"
	WSMsgTypePong         WSMessageType = "pong"
)

// WSMessage estructura genérica de mensaje WebSocket
type WSMessage struct {
	Type      WSMessageType `json:"type"`
	Data      interface{}   `json:"data,omitempty"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// === HELPER FUNCTIONS ===

// CalculateWinPoints calcula los puntos de victoria según racha
func CalculateWinPoints(currentStreak int) int {
	basePoints := 200
	streakBonus := ((currentStreak + 1) / 3) * 100
	return basePoints + streakBonus
}

// CalculateRoundPoints calcula puntos de una ronda
func CalculateRoundPoints(correct bool, timeSeconds float64) int {
	if !correct {
		return 0
	}

	// Puntos base por acierto
	basePoints := 100

	// Bonus por velocidad (máximo 50 puntos)
	// Más rápido = más puntos
	timeBonus := int((15.0 - timeSeconds) / 15.0 * 50.0)
	if timeBonus < 0 {
		timeBonus = 0
	}

	return basePoints + timeBonus
}
