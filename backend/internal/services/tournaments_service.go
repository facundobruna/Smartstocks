package services

import (
	"fmt"
	"time"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type TournamentsService struct {
	tournamentsRepo *repository.TournamentsRepository
	userRepo        *repository.UserRepository
}

func NewTournamentsService(
	tournamentsRepo *repository.TournamentsRepository,
	userRepo *repository.UserRepository,
) *TournamentsService {
	return &TournamentsService{
		tournamentsRepo: tournamentsRepo,
		userRepo:        userRepo,
	}
}

// GetActiveTournaments obtiene torneos activos con detalles
func (s *TournamentsService) GetActiveTournaments(userID string) (*models.TournamentListResponse, error) {
	tournaments, err := s.tournamentsRepo.GetActiveTournaments()
	if err != nil {
		return nil, fmt.Errorf("error getting active tournaments: %w", err)
	}

	// Obtener stats del usuario
	userStats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Enriquecer con detalles
	var enriched []models.TournamentWithDetails
	for _, t := range tournaments {
		details, err := s.enrichTournamentDetails(&t, userID, userStats.RankTier)
		if err != nil {
			continue
		}
		enriched = append(enriched, *details)
	}

	response := &models.TournamentListResponse{
		Tournaments: enriched,
		Total:       len(enriched),
	}

	return response, nil
}

// GetTournamentDetails obtiene detalles de un torneo específico
func (s *TournamentsService) GetTournamentDetails(tournamentID, userID string) (*models.TournamentWithDetails, error) {
	tournament, err := s.tournamentsRepo.GetTournamentByID(tournamentID)
	if err != nil {
		return nil, err
	}

	userStats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, err
	}

	return s.enrichTournamentDetails(tournament, userID, userStats.RankTier)
}

// JoinTournament inscribe al usuario en un torneo
func (s *TournamentsService) JoinTournament(tournamentID, userID string) error {
	success, errorMsg, err := s.tournamentsRepo.JoinTournament(tournamentID, userID)
	if err != nil {
		return fmt.Errorf("error joining tournament: %w", err)
	}

	if !success {
		return fmt.Errorf(errorMsg)
	}

	return nil
}

// GetTournamentStandings obtiene las posiciones del torneo
func (s *TournamentsService) GetTournamentStandings(tournamentID, userID string) (*models.TournamentStandingsResponse, error) {
	tournament, err := s.tournamentsRepo.GetTournamentByID(tournamentID)
	if err != nil {
		return nil, err
	}

	participants, err := s.tournamentsRepo.GetTournamentStandings(tournamentID)
	if err != nil {
		return nil, err
	}

	// Enriquecer con información de usuarios
	var enrichedParticipants []models.TournamentParticipantDetails
	var myPosition *models.TournamentParticipantDetails

	for _, p := range participants {
		user, err := s.userRepo.GetUserByID(p.UserID)
		if err != nil {
			continue
		}

		stats, _ := s.userRepo.GetUserStats(p.UserID)

		detail := models.TournamentParticipantDetails{
			TournamentParticipant: p,
			Username:              user.Username,
			RankTier:              stats.RankTier,
			IsCurrentUser:         p.UserID == userID,
		}

		if user.ProfilePictureURL.Valid {
			detail.ProfilePictureURL = &user.ProfilePictureURL.String
		}

		enrichedParticipants = append(enrichedParticipants, detail)

		if p.UserID == userID {
			detailCopy := detail
			myPosition = &detailCopy
		}
	}

	response := &models.TournamentStandingsResponse{
		TournamentID:   tournamentID,
		TournamentName: tournament.Name,
		Participants:   enrichedParticipants,
		MyPosition:     myPosition,
		TotalPlayers:   len(enrichedParticipants),
	}

	return response, nil
}

// GetTournamentBracket obtiene el bracket del torneo
func (s *TournamentsService) GetTournamentBracket(tournamentID, userID string) (*models.TournamentBracketResponse, error) {
	matches, err := s.tournamentsRepo.GetTournamentMatches(tournamentID)
	if err != nil {
		return nil, err
	}

	// Organizar matches por ronda
	roundsMap := make(map[int][]models.TournamentMatch)
	maxRound := 0

	for _, match := range matches {
		roundsMap[match.RoundNumber] = append(roundsMap[match.RoundNumber], match)
		if match.RoundNumber > maxRound {
			maxRound = match.RoundNumber
		}
	}

	// Crear rounds con información de usuarios
	var rounds []models.TournamentRound
	for roundNum := 1; roundNum <= maxRound; roundNum++ {
		roundMatches := roundsMap[roundNum]
		var enrichedMatches []models.TournamentMatchWithUsers

		for _, match := range roundMatches {
			enriched := s.enrichMatchWithUsers(&match)
			enrichedMatches = append(enrichedMatches, *enriched)
		}

		round := models.TournamentRound{
			RoundNumber: roundNum,
			RoundName:   getRoundName(roundNum, maxRound),
			Matches:     enrichedMatches,
		}

		rounds = append(rounds, round)
	}

	// Determinar ronda actual
	currentRound := 1
	for _, round := range rounds {
		allCompleted := true
		for _, match := range round.Matches {
			if match.Status != "completed" {
				allCompleted = false
				break
			}
		}
		if !allCompleted {
			currentRound = round.RoundNumber
			break
		}
	}

	response := &models.TournamentBracketResponse{
		TournamentID: tournamentID,
		Rounds:       rounds,
		CurrentRound: currentRound,
	}

	return response, nil
}

// GetMyTournaments obtiene los torneos del usuario
func (s *TournamentsService) GetMyTournaments(userID string) (*models.TournamentListResponse, error) {
	tournaments, err := s.tournamentsRepo.GetUserTournaments(userID)
	if err != nil {
		return nil, err
	}

	userStats, _ := s.userRepo.GetUserStats(userID)

	var enriched []models.TournamentWithDetails
	for _, t := range tournaments {
		details, err := s.enrichTournamentDetails(&t, userID, userStats.RankTier)
		if err != nil {
			continue
		}
		enriched = append(enriched, *details)
	}

	response := &models.TournamentListResponse{
		Tournaments: enriched,
		Total:       len(enriched),
	}

	return response, nil
}

// === HELPERS ===

func (s *TournamentsService) enrichTournamentDetails(tournament *models.Tournament, userID, userRank string) (*models.TournamentWithDetails, error) {
	// Obtener premios
	prizes, err := s.tournamentsRepo.GetTournamentPrizes(tournament.ID)
	if err != nil {
		prizes = []models.TournamentPrize{}
	}

	// Verificar si está registrado
	isRegistered, _ := s.tournamentsRepo.IsUserRegistered(tournament.ID, userID)

	// Verificar si puede registrarse
	now := time.Now()
	registrationOpen := tournament.Status == models.TournamentStatusRegistration &&
		now.After(tournament.RegistrationStart) &&
		now.Before(tournament.RegistrationEnd)

	canRegister := !isRegistered &&
		registrationOpen &&
		tournament.CurrentParticipants < tournament.MaxParticipants &&
		meetsRankRequirement(userRank, tournament.MinRankRequired)

	// Calcular tiempo hasta inicio
	var timeUntilStart string
	if tournament.Status == models.TournamentStatusUpcoming || tournament.Status == models.TournamentStatusRegistration {
		duration := time.Until(tournament.StartTime)
		timeUntilStart = formatDuration(duration)
	}

	details := &models.TournamentWithDetails{
		Tournament:       *tournament,
		Prizes:           prizes,
		IsRegistered:     isRegistered,
		CanRegister:      canRegister,
		RegistrationOpen: registrationOpen,
		TimeUntilStart:   timeUntilStart,
		SpotsRemaining:   tournament.MaxParticipants - tournament.CurrentParticipants,
	}

	return details, nil
}

func (s *TournamentsService) enrichMatchWithUsers(match *models.TournamentMatch) *models.TournamentMatchWithUsers {
	enriched := &models.TournamentMatchWithUsers{
		TournamentMatch: *match,
	}

	// Obtener info de player1
	if user1, err := s.userRepo.GetUserByID(match.Player1ID); err == nil {
		enriched.Player1 = &models.UserInfo{
			ID:       user1.ID,
			Username: user1.Username,
		}
		if user1.ProfilePictureURL.Valid {
			enriched.Player1.ProfilePictureURL = &user1.ProfilePictureURL.String
		}
	}

	// Obtener info de player2
	if user2, err := s.userRepo.GetUserByID(match.Player2ID); err == nil {
		enriched.Player2 = &models.UserInfo{
			ID:       user2.ID,
			Username: user2.Username,
		}
		if user2.ProfilePictureURL.Valid {
			enriched.Player2.ProfilePictureURL = &user2.ProfilePictureURL.String
		}
	}

	// Obtener info del ganador
	if match.WinnerID != nil {
		if winner, err := s.userRepo.GetUserByID(*match.WinnerID); err == nil {
			enriched.Winner = &models.UserInfo{
				ID:       winner.ID,
				Username: winner.Username,
			}
			if winner.ProfilePictureURL.Valid {
				enriched.Winner.ProfilePictureURL = &winner.ProfilePictureURL.String
			}
		}
	}

	return enriched
}

func meetsRankRequirement(userRank, requiredRank string) bool {
	rankValues := map[string]int{
		"Bronze 1": 1, "Bronze 2": 2, "Bronze 3": 3,
		"Plata 1": 4, "Plata 2": 5, "Plata 3": 6,
		"Oro 1": 7, "Oro 2": 8, "Oro 3": 9,
		"Maestro": 10,
	}

	userValue := rankValues[userRank]
	requiredValue := rankValues[requiredRank]

	return userValue >= requiredValue
}

func getRoundName(roundNum, maxRound int) string {
	remaining := maxRound - roundNum + 1
	switch remaining {
	case 1:
		return "Final"
	case 2:
		return "Semifinal"
	case 3:
		return "Cuartos de Final"
	default:
		return fmt.Sprintf("Ronda %d", roundNum)
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		return "Started"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
