package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type SimulatorService struct {
	simulatorRepo *repository.SimulatorRepository
	userRepo      *repository.UserRepository
	aiService     *SimulatorAIService
}

func NewSimulatorService(
	simulatorRepo *repository.SimulatorRepository,
	userRepo *repository.UserRepository,
	aiService *SimulatorAIService,
) *SimulatorService {
	return &SimulatorService{
		simulatorRepo: simulatorRepo,
		userRepo:      userRepo,
		aiService:     aiService,
	}
}

// GetScenario obtiene o genera un escenario para el usuario
func (s *SimulatorService) GetScenario(userID string, difficulty models.SimulatorDifficulty) (*models.SimulatorScenarioResponse, error) {
	// Verificar cooldown
	canAttempt, err := s.simulatorRepo.CheckCooldown(userID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("error checking cooldown: %w", err)
	}

	if !canAttempt {
		return nil, errors.New("you have already attempted this difficulty today. Try again tomorrow")
	}

	// Buscar escenario activo existente
	scenario, err := s.simulatorRepo.GetActiveScenarioByDifficulty(difficulty)
	if err != nil {
		return nil, fmt.Errorf("error getting scenario: %w", err)
	}

	// Si no hay escenario, generar uno nuevo
	if scenario == nil {
		scenario, err = s.generateNewScenario(difficulty)
		if err != nil {
			return nil, fmt.Errorf("error generating scenario: %w", err)
		}
	}

	// Preparar respuesta (sin revelar la respuesta correcta)
	response := &models.SimulatorScenarioResponse{
		ScenarioID:  scenario.ID,
		Difficulty:  scenario.Difficulty,
		NewsContent: scenario.NewsContent,
		ChartData: models.ChartData{
			Labels:    scenario.ChartData.Labels,
			Prices:    scenario.ChartData.Prices, // Solo los precios visibles
			Ticker:    scenario.ChartData.Ticker,
			AssetName: scenario.ChartData.AssetName,
			// NO incluir FullPrices ni CorrectDecision
		},
		ExpiresAt: scenario.ExpiresAt,
	}

	return response, nil
}

// SubmitDecision procesa la decisión del usuario
func (s *SimulatorService) SubmitDecision(userID string, req *models.SubmitSimulatorDecisionRequest) (*models.SubmitSimulatorDecisionResponse, error) {
	// Obtener el escenario
	scenario, err := s.simulatorRepo.GetScenarioByID(req.ScenarioID)
	if err != nil {
		return nil, fmt.Errorf("error getting scenario: %w", err)
	}

	// Verificar que el escenario esté activo y no expirado
	if !scenario.IsActive || scenario.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("scenario is no longer active or has expired")
	}

	// Verificar cooldown (por si acaso)
	canAttempt, err := s.simulatorRepo.CheckCooldown(userID, scenario.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("error checking cooldown: %w", err)
	}
	if !canAttempt {
		return nil, errors.New("you have already attempted this difficulty today")
	}

	// Evaluar decisión
	wasCorrect := req.Decision == scenario.CorrectDecision
	pointsEarned := 0
	if wasCorrect {
		pointsEarned = scenario.Difficulty.GetPoints()
	}

	// Registrar intento
	attempt := &models.SimulatorAttempt{
		UserID:       userID,
		ScenarioID:   req.ScenarioID,
		Difficulty:   scenario.Difficulty,
		UserDecision: req.Decision,
		WasCorrect:   wasCorrect,
		PointsEarned: pointsEarned,
		CreatedAt:    time.Now(),
	}

	if req.TimeTakenSeconds != nil {
		attempt.TimeTakenSeconds = sql.NullInt64{
			Int64: int64(*req.TimeTakenSeconds),
			Valid: true,
		}
	}

	// Guardar intento (esto también actualiza stats del usuario)
	if err := s.simulatorRepo.RecordAttempt(attempt); err != nil {
		return nil, fmt.Errorf("error recording attempt: %w", err)
	}

	// Obtener stats actualizados del usuario
	userStats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Preparar respuesta
	response := &models.SubmitSimulatorDecisionResponse{
		WasCorrect:      wasCorrect,
		CorrectDecision: scenario.CorrectDecision,
		UserDecision:    req.Decision,
		PointsEarned:    pointsEarned,
		Explanation:     scenario.Explanation,
		FullChartData:   scenario.ChartData, // Ahora sí mostramos el gráfico completo
		NewTotalPoints:  userStats.Smartpoints,
		NewRankTier:     userStats.RankTier,
	}

	return response, nil
}

// GetHistory obtiene el historial de intentos del usuario
func (s *SimulatorService) GetHistory(userID string, limit int) (*models.SimulatorHistoryResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	attempts, err := s.simulatorRepo.GetUserAttempts(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting attempts: %w", err)
	}

	stats, err := s.simulatorRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %w", err)
	}

	response := &models.SimulatorHistoryResponse{
		Attempts: attempts,
		Stats:    *stats,
	}

	return response, nil
}

// GetCooldownStatus obtiene el estado del cooldown para una dificultad
func (s *SimulatorService) GetCooldownStatus(userID string, difficulty models.SimulatorDifficulty) (*models.CooldownStatusResponse, error) {
	canAttempt, err := s.simulatorRepo.CheckCooldown(userID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("error checking cooldown: %w", err)
	}

	response := &models.CooldownStatusResponse{
		CanAttempt: canAttempt,
	}

	if !canAttempt {
		// Obtener información del último intento
		cooldown, err := s.simulatorRepo.GetLastCooldown(userID, difficulty)
		if err != nil {
			return nil, fmt.Errorf("error getting cooldown: %w", err)
		}

		if cooldown != nil {
			response.LastAttemptDate = &cooldown.LastAttemptDate

			// Calcular próximo disponible (medianoche del día siguiente)
			nextDay := cooldown.LastAttemptDate.AddDate(0, 0, 1)
			nextAvailable := time.Date(
				nextDay.Year(), nextDay.Month(), nextDay.Day(),
				0, 0, 0, 0, nextDay.Location(),
			)
			response.NextAvailable = &nextAvailable

			// Calcular horas restantes
			hoursRemaining := time.Until(nextAvailable).Hours()
			response.HoursRemaining = &hoursRemaining
		}
	}

	return response, nil
}

// generateNewScenario genera un nuevo escenario usando IA
func (s *SimulatorService) generateNewScenario(difficulty models.SimulatorDifficulty) (*models.SimulatorScenario, error) {
	// Intentar generar con IA
	scenario, err := s.aiService.GenerateScenario(difficulty)
	if err != nil {
		// Si falla la IA, usar escenario de respaldo
		fmt.Printf("Warning: AI generation failed, using fallback. Error: %v\n", err)
		scenario = s.aiService.GenerateFallbackScenario(difficulty)
	}

	// Guardar en base de datos
	if err := s.simulatorRepo.CreateScenario(scenario); err != nil {
		return nil, fmt.Errorf("error saving scenario: %w", err)
	}

	return scenario, nil
}

// CleanupExpiredScenarios limpia escenarios expirados (para ejecutar periódicamente)
func (s *SimulatorService) CleanupExpiredScenarios() error {
	return s.simulatorRepo.CleanupExpiredScenarios()
}
