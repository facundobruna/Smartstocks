package services

import (
	"errors"
	"fmt"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type PvPService struct {
	pvpRepo       *repository.PvPRepository
	simulatorRepo *repository.SimulatorRepository
	userRepo      *repository.UserRepository
	aiService     *SimulatorAIService
}

func NewPvPService(
	pvpRepo *repository.PvPRepository,
	simulatorRepo *repository.SimulatorRepository,
	userRepo *repository.UserRepository,
	aiService *SimulatorAIService,
) *PvPService {
	return &PvPService{
		pvpRepo:       pvpRepo,
		simulatorRepo: simulatorRepo,
		userRepo:      userRepo,
		aiService:     aiService,
	}
}

// JoinQueue añade un usuario a la cola de matchmaking
func (s *PvPService) JoinQueue(userID string) (*models.JoinQueueResponse, error) {
	// Obtener stats del usuario para saber su rango
	stats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Unirse a la cola
	entry, err := s.pvpRepo.JoinQueue(userID, stats.RankTier)
	if err != nil {
		return nil, fmt.Errorf("error joining queue: %w", err)
	}

	// Obtener posición en la cola
	position, err := s.pvpRepo.GetQueuePosition(userID)
	if err != nil {
		position = 1 // Por defecto
	}

	response := &models.JoinQueueResponse{
		QueueID:   entry.ID,
		Position:  position,
		ExpiresAt: entry.ExpiresAt,
		Message:   "You have joined the queue. Searching for opponent...",
	}

	return response, nil
}

// LeaveQueue saca un usuario de la cola
func (s *PvPService) LeaveQueue(userID string) error {
	return s.pvpRepo.LeaveQueue(userID)
}

// FindMatch busca un oponente y crea una partida
func (s *PvPService) FindMatch(userID string) (*models.MatchFoundResponse, error) {
	// Obtener stats del usuario
	stats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Buscar oponente
	opponent, err := s.pvpRepo.FindOpponent(userID, stats.RankTier)
	if err != nil {
		return nil, fmt.Errorf("error finding opponent: %w", err)
	}

	if opponent == nil {
		return nil, nil // No hay oponente disponible
	}

	// Crear partida
	match, err := s.pvpRepo.CreateMatch(userID, opponent.UserID)
	if err != nil {
		return nil, fmt.Errorf("error creating match: %w", err)
	}

	// Remover ambos usuarios de la cola
	_ = s.pvpRepo.LeaveQueue(userID)
	_ = s.pvpRepo.LeaveQueue(opponent.UserID)

	// Obtener info del oponente
	opponentUser, err := s.userRepo.GetUserByID(opponent.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting opponent info: %w", err)
	}

	response := &models.MatchFoundResponse{
		MatchID:     match.ID,
		OpponentID:  opponent.UserID,
		Opponent:    s.userToUserInfo(opponentUser),
		TotalRounds: match.TotalRounds,
		Message:     "Match found! Get ready...",
	}

	return response, nil
}

// StartRound inicia una nueva ronda
func (s *PvPService) StartRound(matchID string, roundNumber int) (*models.RoundStartResponse, error) {
	// Verificar que la partida existe
	match, err := s.pvpRepo.GetMatchByID(matchID)
	if err != nil {
		return nil, err
	}

	// Si es la primera ronda, marcar partida como iniciada
	if roundNumber == 1 {
		if err := s.pvpRepo.StartMatch(matchID); err != nil {
			return nil, err
		}
	}

	// Generar escenario para PvP (dificultad medium por defecto)
	scenario, err := s.generatePvPScenario()
	if err != nil {
		return nil, fmt.Errorf("error generating scenario: %w", err)
	}

	// Crear la ronda
	_, err = s.pvpRepo.CreateRound(matchID, roundNumber, scenario.ID, scenario.CorrectDecision)
	if err != nil {
		return nil, fmt.Errorf("error creating round: %w", err)
	}

	// Preparar respuesta (sin revelar respuesta correcta)
	response := &models.RoundStartResponse{
		MatchID:     matchID,
		RoundNumber: roundNumber,
		TotalRounds: match.TotalRounds,
		Scenario: models.SimulatorScenarioResponse{
			ScenarioID:  scenario.ID,
			Difficulty:  scenario.Difficulty,
			NewsContent: scenario.NewsContent,
			ChartData: models.ChartData{
				Labels:    scenario.ChartData.Labels,
				Prices:    scenario.ChartData.Prices,
				Ticker:    scenario.ChartData.Ticker,
				AssetName: scenario.ChartData.AssetName,
			},
			ExpiresAt: scenario.ExpiresAt,
		},
		TimeLimit: 15, // 15 segundos por ronda
	}

	return response, nil
}

// SubmitDecision registra la decisión de un jugador en una ronda
func (s *PvPService) SubmitDecision(userID string, req *models.SubmitPvPDecisionRequest) (*models.RoundResultResponse, error) {
	// Verificar que la partida existe
	match, err := s.pvpRepo.GetMatchByID(req.MatchID)
	if err != nil {
		return nil, err
	}

	// Verificar que el usuario es parte de la partida
	if userID != match.Player1ID && userID != match.Player2ID {
		return nil, errors.New("you are not part of this match")
	}

	// Registrar decisión
	err = s.pvpRepo.SubmitRoundDecision(req.MatchID, req.RoundNumber, userID, req.Decision, req.TimeElapsed)
	if err != nil {
		return nil, err
	}

	// Obtener la ronda actualizada
	round, err := s.pvpRepo.GetRound(req.MatchID, req.RoundNumber)
	if err != nil {
		return nil, err
	}

	// Si ambos jugadores han decidido, completar la ronda
	if round.Player1Decision.Valid && round.Player2Decision.Valid {
		// Completar ronda (calcula puntos)
		if err := s.pvpRepo.CompleteRound(req.MatchID, req.RoundNumber); err != nil {
			return nil, err
		}

		// Recargar ronda con puntos calculados
		round, err = s.pvpRepo.GetRound(req.MatchID, req.RoundNumber)
		if err != nil {
			return nil, err
		}

		// Actualizar puntajes de la partida
		newPlayer1Score := match.Player1Score + round.Player1Points
		newPlayer2Score := match.Player2Score + round.Player2Points
		if err := s.pvpRepo.UpdateMatchScores(req.MatchID, newPlayer1Score, newPlayer2Score); err != nil {
			return nil, err
		}

		// Obtener explicación del escenario
		scenario, err := s.simulatorRepo.GetScenarioByID(round.ScenarioID)
		if err != nil {
			return nil, err
		}

		// Determinar perspectiva del usuario
		isPlayer1 := userID == match.Player1ID

		yourDecision := round.Player1Decision.String
		opponentDecision := round.Player2Decision.String
		yourTime := round.Player1TimeSeconds.Float64
		opponentTime := round.Player2TimeSeconds.Float64
		yourPoints := round.Player1Points
		opponentPoints := round.Player2Points
		yourCorrect := round.Player1Correct.Bool
		opponentCorrect := round.Player2Correct.Bool
		yourTotalScore := newPlayer1Score
		opponentTotalScore := newPlayer2Score

		if !isPlayer1 {
			yourDecision = round.Player2Decision.String
			opponentDecision = round.Player1Decision.String
			yourTime = round.Player2TimeSeconds.Float64
			opponentTime = round.Player1TimeSeconds.Float64
			yourPoints = round.Player2Points
			opponentPoints = round.Player1Points
			yourCorrect = round.Player2Correct.Bool
			opponentCorrect = round.Player1Correct.Bool
			yourTotalScore = newPlayer2Score
			opponentTotalScore = newPlayer1Score
		}

		// Verificar si la partida está completa
		isComplete := req.RoundNumber >= match.TotalRounds

		response := &models.RoundResultResponse{
			MatchID:            req.MatchID,
			RoundNumber:        req.RoundNumber,
			YourDecision:       models.SimulatorDecision(yourDecision),
			OpponentDecision:   models.SimulatorDecision(opponentDecision),
			CorrectDecision:    round.CorrectDecision,
			YourCorrect:        yourCorrect,
			OpponentCorrect:    opponentCorrect,
			YourTime:           yourTime,
			OpponentTime:       opponentTime,
			YourPoints:         yourPoints,
			OpponentPoints:     opponentPoints,
			YourTotalScore:     yourTotalScore,
			OpponentTotalScore: opponentTotalScore,
			Explanation:        scenario.Explanation,
			IsMatchComplete:    isComplete,
		}

		return response, nil
	}

	// Si solo ha decidido un jugador, esperar al otro
	return nil, errors.New("waiting for opponent decision")
}

// GetMatchResult obtiene el resultado final de una partida
func (s *PvPService) GetMatchResult(userID, matchID string) (*models.MatchResultResponse, error) {
	// Obtener partida
	match, err := s.pvpRepo.GetMatchByID(matchID)
	if err != nil {
		return nil, err
	}

	// Verificar que el usuario es parte de la partida
	if userID != match.Player1ID && userID != match.Player2ID {
		return nil, errors.New("you are not part of this match")
	}

	// Determinar ganador
	var winner string
	var winnerID string
	var loserID string
	var pointsGained int

	isPlayer1 := userID == match.Player1ID
	yourScore := match.Player1Score
	opponentScore := match.Player2Score
	if !isPlayer1 {
		yourScore = match.Player2Score
		opponentScore = match.Player1Score
	}

	if yourScore > opponentScore {
		winner = "you"
		winnerID = userID
		if isPlayer1 {
			loserID = match.Player2ID
		} else {
			loserID = match.Player1ID
		}
	} else if opponentScore > yourScore {
		winner = "opponent"
		if isPlayer1 {
			winnerID = match.Player2ID
			loserID = match.Player1ID
		} else {
			winnerID = match.Player1ID
			loserID = match.Player2ID
		}
	} else {
		winner = "tie"
	}

	// Calcular puntos ganados/perdidos
	userStats, _ := s.userRepo.GetUserStats(userID)
	streakBonus := 0

	if winner == "you" {
		pointsGained = models.CalculateWinPoints(userStats.WinStreak)
		streakBonus = pointsGained - 200 // Base es 200

		// Actualizar stats
		_ = s.pvpRepo.UpdatePvPStats(winnerID, loserID, pointsGained)
	} else if winner == "opponent" {
		pointsGained = -100 // Pierde 100 puntos

		// Actualizar stats
		_ = s.pvpRepo.UpdatePvPStats(winnerID, loserID, 200) // El ganador recibe base
	} else {
		pointsGained = 0 // Empate, no gana ni pierde
	}

	// Marcar partida como completada
	if winner != "tie" {
		_ = s.pvpRepo.CompleteMatch(matchID, winnerID)
	} else {
		_ = s.pvpRepo.CompleteMatch(matchID, "") // Sin ganador
	}

	// Obtener stats actualizados
	newStats, _ := s.userRepo.GetUserStats(userID)

	// Obtener resumen de rondas
	rounds, err := s.pvpRepo.GetMatchRounds(matchID)
	if err != nil {
		return nil, err
	}

	roundSummaries := make([]models.RoundSummary, len(rounds))
	for i, round := range rounds {
		yourDec := round.Player1Decision.String
		oppDec := round.Player2Decision.String
		yourPts := round.Player1Points
		oppPts := round.Player2Points

		if !isPlayer1 {
			yourDec = round.Player2Decision.String
			oppDec = round.Player1Decision.String
			yourPts = round.Player2Points
			oppPts = round.Player1Points
		}

		roundSummaries[i] = models.RoundSummary{
			RoundNumber:      round.RoundNumber,
			YourDecision:     models.SimulatorDecision(yourDec),
			OpponentDecision: models.SimulatorDecision(oppDec),
			CorrectDecision:  round.CorrectDecision,
			YourPoints:       yourPts,
			OpponentPoints:   oppPts,
		}
	}

	response := &models.MatchResultResponse{
		MatchID:            matchID,
		Winner:             winner,
		YourFinalScore:     yourScore,
		OpponentFinalScore: opponentScore,
		PointsGained:       pointsGained,
		NewTotalPoints:     newStats.Smartpoints,
		NewRankTier:        newStats.RankTier,
		WinStreak:          newStats.WinStreak,
		StreakBonus:        streakBonus,
		Rounds:             roundSummaries,
	}

	return response, nil
}

// GetHistory obtiene el historial de partidas
func (s *PvPService) GetHistory(userID string, limit int) (*models.PvPHistoryResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	matches, err := s.pvpRepo.GetUserMatches(userID, limit)
	if err != nil {
		return nil, err
	}

	// Enriquecer con información de usuarios
	matchesWithDetails := make([]models.PvPMatchWithDetails, len(matches))
	for i, match := range matches {
		player1, _ := s.userRepo.GetUserByID(match.Player1ID)
		player2, _ := s.userRepo.GetUserByID(match.Player2ID)

		matchDetail := models.PvPMatchWithDetails{
			PvPMatch: match,
			Player1:  s.userToUserInfo(player1),
			Player2:  s.userToUserInfo(player2),
		}

		if match.WinnerID.Valid {
			winner, _ := s.userRepo.GetUserByID(match.WinnerID.String)
			matchDetail.Winner = s.userToUserInfo(winner)
		}

		matchesWithDetails[i] = matchDetail
	}

	// Obtener stats
	stats, err := s.pvpRepo.GetUserPvPStats(userID)
	if err != nil {
		return nil, err
	}

	response := &models.PvPHistoryResponse{
		Matches: matchesWithDetails,
		Stats:   *stats,
	}

	return response, nil
}

// === HELPERS ===

func (s *PvPService) generatePvPScenario() (*models.SimulatorScenario, error) {
	// Buscar escenario aleatorio de dificultad medium
	scenario, err := s.simulatorRepo.GetRandomScenarioByDifficulty(models.SimulatorDifficultyMedium)
	if err != nil || scenario == nil {
		// Si no hay escenarios, generar uno nuevo
		scenario, err = s.aiService.GenerateScenario(models.SimulatorDifficultyMedium)
		if err != nil {
			// Usar fallback
			scenario = s.aiService.GenerateFallbackScenario(models.SimulatorDifficultyMedium)
		}

		// Guardar
		if err := s.simulatorRepo.CreateScenario(scenario); err != nil {
			return nil, err
		}
	}

	return scenario, nil
}

func (s *PvPService) userToUserInfo(user *models.User) *models.UserInfo {
	if user == nil {
		return nil
	}

	userInfo := &models.UserInfo{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.ProfilePictureURL.Valid {
		userInfo.ProfilePictureURL = &user.ProfilePictureURL.String
	}

	if user.SchoolID.Valid {
		userInfo.SchoolID = &user.SchoolID.String
	}

	return userInfo
}
