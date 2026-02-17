package models

import (
	"time"
)

// LeaderboardEntry representa una entrada en el ranking
type LeaderboardEntry struct {
	RankPosition      int     `json:"rank_position"`
	UserID            string  `json:"user_id"`
	Username          string  `json:"username"`
	Smartpoints       int     `json:"smartpoints"`
	RankTier          string  `json:"rank_tier"`
	TotalWins         int     `json:"total_wins"`
	TotalLosses       int     `json:"total_losses"`
	WinRate           float64 `json:"win_rate"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	SchoolName        *string `json:"school_name,omitempty"`
	SchoolID          *string `json:"school_id,omitempty"`
	IsCurrentUser     bool    `json:"is_current_user"`
}

// LeaderboardResponse representa la respuesta del ranking
type LeaderboardResponse struct {
	Type         string             `json:"type"` // "global" o "school"
	TopPlayers   []LeaderboardEntry `json:"top_players"`
	UserPosition *LeaderboardEntry  `json:"user_position,omitempty"`
	TotalPlayers int                `json:"total_players"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// UserPositionResponse representa la posición del usuario
type UserPositionResponse struct {
	GlobalPosition int `json:"global_position"`
	SchoolPosition int `json:"school_position,omitempty"`
	TotalPlayers   int `json:"total_players"`
}

// Achievement representa un logro
type Achievement struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	AchievementType        string    `json:"achievement_type"`
	AchievementName        string    `json:"achievement_name"`
	AchievementDescription string    `json:"achievement_description"`
	IconURL                *string   `json:"icon_url,omitempty"`
	UnlockedAt             time.Time `json:"unlocked_at"`
	IsUnlocked             bool      `json:"is_unlocked"`
}

// UserProfilePublic representa un perfil público de usuario
type UserProfilePublic struct {
	UserInfo
	Stats        *UserStats    `json:"stats"`
	Achievements []Achievement `json:"achievements"`
	GlobalRank   int           `json:"global_rank"`
	SchoolRank   int           `json:"school_rank,omitempty"`
}

// SchoolLeaderboardResponse representa el ranking de un colegio específico
type SchoolLeaderboardResponse struct {
	SchoolID     string             `json:"school_id"`
	SchoolName   string             `json:"school_name"`
	TopPlayers   []LeaderboardEntry `json:"top_players"`
	TotalPlayers int                `json:"total_players"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// LeaderboardFilters representa filtros para el ranking
type LeaderboardFilters struct {
	Limit      int    `json:"limit" form:"limit"`
	Offset     int    `json:"offset" form:"offset"`
	SchoolID   string `json:"school_id" form:"school_id"`
	RankTier   string `json:"rank_tier" form:"rank_tier"`
	SearchTerm string `json:"search" form:"search"`
}

// AchievementProgress representa el progreso hacia un logro
type AchievementProgress struct {
	AchievementType string  `json:"achievement_type"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Current         int     `json:"current"`
	Required        int     `json:"required"`
	Progress        float64 `json:"progress"` // Porcentaje (0-100)
	IsUnlocked      bool    `json:"is_unlocked"`
}

// AllAchievementsResponse representa todos los logros disponibles
type AllAchievementsResponse struct {
	Unlocked   []Achievement         `json:"unlocked"`
	Locked     []AchievementProgress `json:"locked"`
	TotalCount int                   `json:"total_count"`
}
