package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/smartstocks/backend/internal/models"
)

type RankingsRepository struct {
	db *sql.DB
}

func NewRankingsRepository(db *sql.DB) *RankingsRepository {
	return &RankingsRepository{db: db}
}

// GetGlobalLeaderboard obtiene el ranking global
func (r *RankingsRepository) GetGlobalLeaderboard(limit, offset int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	query := `
		SELECT 
			rank_position, user_id, username, smartpoints, rank_tier,
			total_wins, total_losses, win_rate, profile_picture_url, school_name
		FROM leaderboard_cache
		WHERE cache_type = 'global'
		ORDER BY rank_position ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		var profilePic, schoolName sql.NullString

		err := rows.Scan(
			&entry.RankPosition,
			&entry.UserID,
			&entry.Username,
			&entry.Smartpoints,
			&entry.RankTier,
			&entry.TotalWins,
			&entry.TotalLosses,
			&entry.WinRate,
			&profilePic,
			&schoolName,
		)
		if err != nil {
			return nil, err
		}

		if profilePic.Valid {
			entry.ProfilePictureURL = &profilePic.String
		}
		if schoolName.Valid {
			entry.SchoolName = &schoolName.String
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetSchoolLeaderboard obtiene el ranking de un colegio
func (r *RankingsRepository) GetSchoolLeaderboard(schoolID string, limit, offset int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	query := `
		SELECT 
			rank_position, user_id, username, smartpoints, rank_tier,
			total_wins, total_losses, win_rate, profile_picture_url, school_name
		FROM leaderboard_cache
		WHERE cache_type = 'school' AND school_id = ?
		ORDER BY rank_position ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, schoolID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		var profilePic, schoolName sql.NullString

		err := rows.Scan(
			&entry.RankPosition,
			&entry.UserID,
			&entry.Username,
			&entry.Smartpoints,
			&entry.RankTier,
			&entry.TotalWins,
			&entry.TotalLosses,
			&entry.WinRate,
			&profilePic,
			&schoolName,
		)
		if err != nil {
			return nil, err
		}

		if profilePic.Valid {
			entry.ProfilePictureURL = &profilePic.String
		}
		if schoolName.Valid {
			entry.SchoolName = &schoolName.String
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetUserPosition obtiene la posición del usuario en los rankings
func (r *RankingsRepository) GetUserPosition(userID string) (globalPos, schoolPos int, err error) {
	query := `CALL get_user_position(?, @global_pos, @school_pos)`
	_, err = r.db.Exec(query, userID)
	if err != nil {
		return 0, 0, err
	}

	err = r.db.QueryRow(`SELECT @global_pos, @school_pos`).Scan(&globalPos, &schoolPos)
	return globalPos, schoolPos, err
}

// GetUserRankingEntry obtiene la entrada del usuario en el ranking
func (r *RankingsRepository) GetUserRankingEntry(userID string, leaderboardType string) (*models.LeaderboardEntry, error) {
	query := `
		SELECT 
			rank_position, user_id, username, smartpoints, rank_tier,
			total_wins, total_losses, win_rate, profile_picture_url, school_name
		FROM leaderboard_cache
		WHERE cache_type = ? AND user_id = ?
		LIMIT 1
	`

	var entry models.LeaderboardEntry
	var profilePic, schoolName sql.NullString

	err := r.db.QueryRow(query, leaderboardType, userID).Scan(
		&entry.RankPosition,
		&entry.UserID,
		&entry.Username,
		&entry.Smartpoints,
		&entry.RankTier,
		&entry.TotalWins,
		&entry.TotalLosses,
		&entry.WinRate,
		&profilePic,
		&schoolName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if profilePic.Valid {
		entry.ProfilePictureURL = &profilePic.String
	}
	if schoolName.Valid {
		entry.SchoolName = &schoolName.String
	}

	return &entry, nil
}

// GetTotalPlayers obtiene el total de jugadores
func (r *RankingsRepository) GetTotalPlayers(leaderboardType, schoolID string) (int, error) {
	var count int
	var query string

	if leaderboardType == "school" && schoolID != "" {
		query = `SELECT COUNT(*) FROM leaderboard_cache WHERE cache_type = 'school' AND school_id = ?`
		err := r.db.QueryRow(query, schoolID).Scan(&count)
		return count, err
	}

	query = `SELECT COUNT(*) FROM leaderboard_cache WHERE cache_type = 'global'`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetLastUpdated obtiene la fecha de última actualización del cache
func (r *RankingsRepository) GetLastUpdated() (time.Time, error) {
	var lastUpdated time.Time
	query := `SELECT MAX(last_updated) FROM leaderboard_cache`
	err := r.db.QueryRow(query).Scan(&lastUpdated)
	return lastUpdated, err
}

// UpdateLeaderboardCache actualiza el cache de rankings
func (r *RankingsRepository) UpdateLeaderboardCache() error {
	query := `CALL update_leaderboard_cache()`
	_, err := r.db.Exec(query)
	return err
}

// === ACHIEVEMENTS ===

// GetUserAchievements obtiene los logros de un usuario
func (r *RankingsRepository) GetUserAchievements(userID string) ([]models.Achievement, error) {
	query := `
		SELECT id, user_id, achievement_type, achievement_name, 
			   achievement_description, icon_url, unlocked_at
		FROM user_achievements
		WHERE user_id = ?
		ORDER BY unlocked_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		var iconURL sql.NullString

		err := rows.Scan(
			&achievement.ID,
			&achievement.UserID,
			&achievement.AchievementType,
			&achievement.AchievementName,
			&achievement.AchievementDescription,
			&iconURL,
			&achievement.UnlockedAt,
		)
		if err != nil {
			return nil, err
		}

		if iconURL.Valid {
			achievement.IconURL = &iconURL.String
		}
		achievement.IsUnlocked = true

		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

// GrantAchievement otorga un logro a un usuario
func (r *RankingsRepository) GrantAchievement(userID, achievementType, name, description string) error {
	query := `CALL grant_achievement(?, ?, ?, ?)`
	_, err := r.db.Exec(query, userID, achievementType, name, description)
	return err
}

// HasAchievement verifica si el usuario tiene un logro
func (r *RankingsRepository) HasAchievement(userID, achievementType string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_achievements WHERE user_id = ? AND achievement_type = ?`
	err := r.db.QueryRow(query, userID, achievementType).Scan(&count)
	return count > 0, err
}

// GetPublicProfile obtiene el perfil público de un usuario
func (r *RankingsRepository) GetPublicProfile(userID string) (*models.UserProfilePublic, error) {
	// Este método se implementará en el servicio combinando múltiples queries
	return nil, fmt.Errorf("use service layer for GetPublicProfile")
}
