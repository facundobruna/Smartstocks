package models

import (
	"database/sql"
	"time"
)

// SimulatorDifficulty representa los niveles de dificultad
type SimulatorDifficulty string

const (
	SimulatorDifficultyEasy   SimulatorDifficulty = "easy"
	SimulatorDifficultyMedium SimulatorDifficulty = "medium"
	SimulatorDifficultyHard   SimulatorDifficulty = "hard"
)

// SimulatorDecision representa las decisiones posibles
type SimulatorDecision string

const (
	SimulatorDecisionBuy  SimulatorDecision = "buy"
	SimulatorDecisionSell SimulatorDecision = "sell"
	SimulatorDecisionHold SimulatorDecision = "hold"
)

// ChartData representa los datos del gráfico en formato JSON
type ChartData struct {
	Labels     []string  `json:"labels"`      // Fechas/tiempo
	Prices     []float64 `json:"prices"`      // Precios hasta el punto de decisión
	FullPrices []float64 `json:"full_prices"` // Precios completos (incluye futuro)
	Ticker     string    `json:"ticker"`      // Símbolo del activo (ej: "AAPL", "BTC")
	AssetName  string    `json:"asset_name"`  // Nombre del activo
}

// SimulatorScenario representa un escenario generado por IA
type SimulatorScenario struct {
	ID              string              `json:"id"`
	Difficulty      SimulatorDifficulty `json:"difficulty"`
	NewsContent     string              `json:"news_content"`
	ChartData       ChartData           `json:"chart_data"`
	CorrectDecision SimulatorDecision   `json:"-"` // No enviar al cliente
	Explanation     string              `json:"-"` // No enviar al cliente hasta que responda
	CreatedAt       time.Time           `json:"created_at"`
	ExpiresAt       time.Time           `json:"expires_at"`
	IsActive        bool                `json:"is_active"`
}

// SimulatorAttempt representa un intento de usuario
type SimulatorAttempt struct {
	ID               string              `json:"id"`
	UserID           string              `json:"user_id"`
	ScenarioID       string              `json:"scenario_id"`
	Difficulty       SimulatorDifficulty `json:"difficulty"`
	UserDecision     SimulatorDecision   `json:"user_decision"`
	WasCorrect       bool                `json:"was_correct"`
	PointsEarned     int                 `json:"points_earned"`
	TimeTakenSeconds sql.NullInt64       `json:"time_taken_seconds,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
}

// DailySimulatorCooldown representa el cooldown diario
type DailySimulatorCooldown struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	Difficulty      SimulatorDifficulty `json:"difficulty"`
	LastAttemptDate time.Time           `json:"last_attempt_date"`
	AttemptsCount   int                 `json:"attempts_count"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// === REQUEST/RESPONSE MODELS ===

// GetSimulatorScenarioRequest representa la solicitud para obtener un escenario
type GetSimulatorScenarioRequest struct {
	Difficulty SimulatorDifficulty `json:"difficulty" binding:"required,oneof=easy medium hard"`
}

// SimulatorScenarioResponse representa la respuesta con el escenario (sin revelar respuesta)
type SimulatorScenarioResponse struct {
	ScenarioID  string              `json:"scenario_id"`
	Difficulty  SimulatorDifficulty `json:"difficulty"`
	NewsContent string              `json:"news_content"`
	ChartData   ChartData           `json:"chart_data"` // Solo hasta el punto de decisión
	ExpiresAt   time.Time           `json:"expires_at"`
}

// SubmitSimulatorDecisionRequest representa el envío de una decisión
type SubmitSimulatorDecisionRequest struct {
	ScenarioID       string            `json:"scenario_id" binding:"required"`
	Decision         SimulatorDecision `json:"decision" binding:"required,oneof=buy sell hold"`
	TimeTakenSeconds *int              `json:"time_taken_seconds,omitempty"`
}

// SubmitSimulatorDecisionResponse representa el resultado después de enviar decisión
type SubmitSimulatorDecisionResponse struct {
	WasCorrect      bool              `json:"was_correct"`
	CorrectDecision SimulatorDecision `json:"correct_decision"`
	UserDecision    SimulatorDecision `json:"user_decision"`
	PointsEarned    int               `json:"points_earned"`
	Explanation     string            `json:"explanation"`
	FullChartData   ChartData         `json:"full_chart_data"` // Gráfico completo con el futuro
	NewTotalPoints  int               `json:"new_total_points"`
	NewRankTier     string            `json:"new_rank_tier"`
}

// SimulatorHistoryResponse representa el historial de intentos
type SimulatorHistoryResponse struct {
	Attempts []SimulatorAttemptWithDetails `json:"attempts"`
	Stats    SimulatorStats                `json:"stats"`
}

// SimulatorAttemptWithDetails incluye detalles del escenario
type SimulatorAttemptWithDetails struct {
	SimulatorAttempt
	NewsContent string `json:"news_content"`
	Explanation string `json:"explanation"`
}

// SimulatorStats representa estadísticas del simulador
type SimulatorStats struct {
	TotalAttempts   int                                 `json:"total_attempts"`
	CorrectAttempts int                                 `json:"correct_attempts"`
	AccuracyRate    float64                             `json:"accuracy_rate"`
	TotalPoints     int                                 `json:"total_points"`
	ByDifficulty    map[string]SimulatorDifficultyStats `json:"by_difficulty"`
}

// SimulatorDifficultyStats representa estadísticas por dificultad
type SimulatorDifficultyStats struct {
	Attempts     int     `json:"attempts"`
	Correct      int     `json:"correct"`
	AccuracyRate float64 `json:"accuracy_rate"`
	PointsEarned int     `json:"points_earned"`
}

// CooldownStatusResponse representa el estado del cooldown
type CooldownStatusResponse struct {
	CanAttempt      bool       `json:"can_attempt"`
	LastAttemptDate *time.Time `json:"last_attempt_date,omitempty"`
	NextAvailable   *time.Time `json:"next_available,omitempty"`
	HoursRemaining  *float64   `json:"hours_remaining,omitempty"`
}

// === HELPER FUNCTIONS ===

// GetPoints retorna los puntos según dificultad
func (d SimulatorDifficulty) GetPoints() int {
	switch d {
	case SimulatorDifficultyEasy:
		return 25
	case SimulatorDifficultyMedium:
		return 50
	case SimulatorDifficultyHard:
		return 100
	default:
		return 0
	}
}

// IsValid verifica si la dificultad es válida
func (d SimulatorDifficulty) IsValid() bool {
	switch d {
	case SimulatorDifficultyEasy, SimulatorDifficultyMedium, SimulatorDifficultyHard:
		return true
	default:
		return false
	}
}

// IsValid verifica si la decisión es válida
func (d SimulatorDecision) IsValid() bool {
	switch d {
	case SimulatorDecisionBuy, SimulatorDecisionSell, SimulatorDecisionHold:
		return true
	default:
		return false
	}
}

// String convierte la decisión a string en español
func (d SimulatorDecision) String() string {
	switch d {
	case SimulatorDecisionBuy:
		return "Comprar"
	case SimulatorDecisionSell:
		return "Vender"
	case SimulatorDecisionHold:
		return "Mantener"
	default:
		return "Desconocido"
	}
}
