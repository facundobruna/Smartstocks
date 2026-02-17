package services

import (
	"fmt"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type RankingsService struct {
	rankingsRepo *repository.RankingsRepository
	userRepo     *repository.UserRepository
}

func NewRankingsService(
	rankingsRepo *repository.RankingsRepository,
	userRepo *repository.UserRepository,
) *RankingsService {
	return &RankingsService{
		rankingsRepo: rankingsRepo,
		userRepo:     userRepo,
	}
}

// GetGlobalLeaderboard obtiene el ranking global
func (s *RankingsService) GetGlobalLeaderboard(userID string, limit, offset int) (*models.LeaderboardResponse, error) {
	// Obtener top players
	topPlayers, err := s.rankingsRepo.GetGlobalLeaderboard(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting global leaderboard: %w", err)
	}

	// Marcar al usuario actual
	for i := range topPlayers {
		if topPlayers[i].UserID == userID {
			topPlayers[i].IsCurrentUser = true
		}
	}

	// Obtener posición del usuario (si no está en el top)
	var userPosition *models.LeaderboardEntry
	userInTop := false
	for _, player := range topPlayers {
		if player.UserID == userID {
			userInTop = true
			break
		}
	}

	if !userInTop {
		userEntry, err := s.rankingsRepo.GetUserRankingEntry(userID, "global")
		if err == nil && userEntry != nil {
			userEntry.IsCurrentUser = true
			userPosition = userEntry
		}
	}

	// Obtener total de jugadores
	totalPlayers, err := s.rankingsRepo.GetTotalPlayers("global", "")
	if err != nil {
		totalPlayers = 0
	}

	// Obtener última actualización
	lastUpdated, err := s.rankingsRepo.GetLastUpdated()
	if err != nil {
		lastUpdated = lastUpdated // Mantener zero value
	}

	response := &models.LeaderboardResponse{
		Type:         "global",
		TopPlayers:   topPlayers,
		UserPosition: userPosition,
		TotalPlayers: totalPlayers,
		LastUpdated:  lastUpdated,
	}

	return response, nil
}

// GetSchoolLeaderboard obtiene el ranking de un colegio
func (s *RankingsService) GetSchoolLeaderboard(userID, schoolID string, limit, offset int) (*models.LeaderboardResponse, error) {
	// Obtener top players del colegio
	topPlayers, err := s.rankingsRepo.GetSchoolLeaderboard(schoolID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting school leaderboard: %w", err)
	}

	// Marcar al usuario actual
	for i := range topPlayers {
		if topPlayers[i].UserID == userID {
			topPlayers[i].IsCurrentUser = true
		}
	}

	// Obtener posición del usuario (si no está en el top)
	var userPosition *models.LeaderboardEntry
	userInTop := false
	for _, player := range topPlayers {
		if player.UserID == userID {
			userInTop = true
			break
		}
	}

	if !userInTop {
		userEntry, err := s.rankingsRepo.GetUserRankingEntry(userID, "school")
		if err == nil && userEntry != nil {
			userEntry.IsCurrentUser = true
			userPosition = userEntry
		}
	}

	// Obtener total de jugadores en el colegio
	totalPlayers, err := s.rankingsRepo.GetTotalPlayers("school", schoolID)
	if err != nil {
		totalPlayers = 0
	}

	// Obtener última actualización
	lastUpdated, err := s.rankingsRepo.GetLastUpdated()
	if err != nil {
		lastUpdated = lastUpdated
	}

	response := &models.LeaderboardResponse{
		Type:         "school",
		TopPlayers:   topPlayers,
		UserPosition: userPosition,
		TotalPlayers: totalPlayers,
		LastUpdated:  lastUpdated,
	}

	return response, nil
}

// GetMySchoolLeaderboard obtiene el ranking del colegio del usuario
func (s *RankingsService) GetMySchoolLeaderboard(userID string, limit, offset int) (*models.LeaderboardResponse, error) {
	// Obtener usuario para saber su colegio
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	if !user.SchoolID.Valid {
		return nil, fmt.Errorf("user does not belong to any school")
	}

	return s.GetSchoolLeaderboard(userID, user.SchoolID.String, limit, offset)
}

// GetUserPosition obtiene la posición del usuario en los rankings
func (s *RankingsService) GetUserPosition(userID string) (*models.UserPositionResponse, error) {
	globalPos, schoolPos, err := s.rankingsRepo.GetUserPosition(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user position: %w", err)
	}

	totalPlayers, _ := s.rankingsRepo.GetTotalPlayers("global", "")

	response := &models.UserPositionResponse{
		GlobalPosition: globalPos,
		SchoolPosition: schoolPos,
		TotalPlayers:   totalPlayers,
	}

	return response, nil
}

// GetPublicProfile obtiene el perfil público de un usuario
func (s *RankingsService) GetPublicProfile(targetUserID, requestingUserID string) (*models.UserProfilePublic, error) {
	// Obtener usuario
	user, err := s.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Obtener stats
	stats, err := s.userRepo.GetUserStats(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Obtener logros
	achievements, err := s.rankingsRepo.GetUserAchievements(targetUserID)
	if err != nil {
		achievements = []models.Achievement{} // No fallar si no hay logros
	}

	// Obtener posiciones
	globalPos, schoolPos, _ := s.rankingsRepo.GetUserPosition(targetUserID)

	// Construir perfil público
	profile := &models.UserProfilePublic{
		UserInfo: models.UserInfo{
			ID:            user.ID,
			Username:      user.Username,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Stats:        stats,
		Achievements: achievements,
		GlobalRank:   globalPos,
		SchoolRank:   schoolPos,
	}

	// No mostrar email a otros usuarios
	if targetUserID != requestingUserID {
		profile.Email = ""
	} else {
		profile.Email = user.Email
	}

	if user.ProfilePictureURL.Valid {
		profile.ProfilePictureURL = &user.ProfilePictureURL.String
	}

	if user.SchoolID.Valid {
		profile.SchoolID = &user.SchoolID.String
	}

	return profile, nil
}

// GetUserAchievements obtiene los logros de un usuario
func (s *RankingsService) GetUserAchievements(userID string) (*models.AllAchievementsResponse, error) {
	// Obtener logros desbloqueados
	unlocked, err := s.rankingsRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting achievements: %w", err)
	}

	// Obtener stats para calcular progreso
	stats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	// Definir todos los logros posibles y calcular progreso
	locked := s.calculateAchievementProgress(userID, stats, unlocked)

	response := &models.AllAchievementsResponse{
		Unlocked:   unlocked,
		Locked:     locked,
		TotalCount: len(unlocked) + len(locked),
	}

	return response, nil
}

// calculateAchievementProgress calcula el progreso de logros bloqueados
func (s *RankingsService) calculateAchievementProgress(userID string, stats *models.UserStats, unlocked []models.Achievement) []models.AchievementProgress {
	// Crear mapa de logros desbloqueados
	unlockedMap := make(map[string]bool)
	for _, achievement := range unlocked {
		unlockedMap[achievement.AchievementType] = true
	}

	var locked []models.AchievementProgress

	// Definir todos los logros posibles
	allAchievements := []struct {
		Type        string
		Name        string
		Description string
		Required    int
		Current     int
	}{
		{"first_win", "Primera Victoria", "Gana tu primera partida PvP", 1, stats.TotalWins},
		{"win_streak_3", "En Racha", "Consigue 3 victorias seguidas", 3, stats.WinStreak},
		{"win_streak_5", "Imparable", "Consigue 5 victorias seguidas", 5, stats.WinStreak},
		{"win_streak_10", "Leyenda", "Consigue 10 victorias seguidas", 10, stats.WinStreak},
		{"rank_bronze", "Rango Bronce", "Alcanza el rango Bronce", 1, boolToInt(stats.RankTier != "")},
		{"rank_silver", "Ascenso a Plata", "Alcanza el rango Plata", 1, boolToInt(containsRank(stats.RankTier, "Plata", "Oro", "Maestro"))},
		{"rank_gold", "Ascenso a Oro", "Alcanza el rango Oro", 1, boolToInt(containsRank(stats.RankTier, "Oro", "Maestro"))},
		{"rank_master", "Maestro de las Finanzas", "Alcanza el rango Maestro", 1, boolToInt(stats.RankTier == "Maestro")},
		{"quiz_master", "Maestro de Quizzes", "Completa 50 quizzes", 50, stats.TotalQuizzesCompleted},
		{"simulator_expert", "Experto Simulador", "Completa 100 simulaciones", 100, stats.TotalSimulatorGames},
		{"pvp_legend", "Leyenda PvP", "Gana 100 partidas PvP", 100, stats.TotalWins},
		{"points_1000", "Mil Puntos", "Alcanza 1,000 SmartPoints", 1000, stats.Smartpoints},
		{"points_5000", "Cinco Mil", "Alcanza 5,000 SmartPoints", 5000, stats.Smartpoints},
		{"points_10000", "Diez Mil", "Alcanza 10,000 SmartPoints", 10000, stats.Smartpoints},
	}

	for _, ach := range allAchievements {
		// Si ya está desbloqueado, saltar
		if unlockedMap[ach.Type] {
			continue
		}

		// Calcular progreso
		progress := float64(0)
		if ach.Required > 0 {
			progress = (float64(ach.Current) / float64(ach.Required)) * 100
			if progress > 100 {
				progress = 100
			}
		}

		locked = append(locked, models.AchievementProgress{
			AchievementType: ach.Type,
			Name:            ach.Name,
			Description:     ach.Description,
			Current:         ach.Current,
			Required:        ach.Required,
			Progress:        progress,
			IsUnlocked:      false,
		})
	}

	return locked
}

// UpdateLeaderboardCache fuerza actualización del cache
func (s *RankingsService) UpdateLeaderboardCache() error {
	return s.rankingsRepo.UpdateLeaderboardCache()
}

// === HELPERS ===

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func containsRank(current string, ranks ...string) bool {
	for _, rank := range ranks {
		if len(current) >= len(rank) && current[:len(rank)] == rank {
			return true
		}
	}
	return false
}
