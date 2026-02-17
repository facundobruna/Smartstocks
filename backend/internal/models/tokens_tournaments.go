package models

import (
	"database/sql"
	"time"
)

// === TOKENS ===

// UserTokens representa el balance de tokens de un usuario
type UserTokens struct {
	UserID            string       `json:"user_id"`
	Balance           int          `json:"balance"`
	TotalEarned       int          `json:"total_earned"`
	TotalSpent        int          `json:"total_spent"`
	LastTransactionAt sql.NullTime `json:"last_transaction_at,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// TokenTransaction representa una transacción de tokens
type TokenTransaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          int       `json:"amount"`
	BalanceAfter    int       `json:"balance_after"`
	Description     string    `json:"description"`
	ReferenceID     *string   `json:"reference_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// TokensResponse representa la respuesta de tokens
type TokensResponse struct {
	Balance            int                `json:"balance"`
	TotalEarned        int                `json:"total_earned"`
	TotalSpent         int                `json:"total_spent"`
	RecentTransactions []TokenTransaction `json:"recent_transactions"`
}

// === TOURNAMENTS ===

// TournamentType tipos de torneo
type TournamentType string

const (
	TournamentTypeWeekly  TournamentType = "weekly"
	TournamentTypeMonthly TournamentType = "monthly"
	TournamentTypeSpecial TournamentType = "special"
)

// TournamentFormat formatos de torneo
type TournamentFormat string

const (
	TournamentFormatBracket      TournamentFormat = "bracket"
	TournamentFormatLeague       TournamentFormat = "league"
	TournamentFormatBattleRoyale TournamentFormat = "battle_royale"
)

// TournamentStatus estados del torneo
type TournamentStatus string

const (
	TournamentStatusUpcoming     TournamentStatus = "upcoming"
	TournamentStatusRegistration TournamentStatus = "registration"
	TournamentStatusInProgress   TournamentStatus = "in_progress"
	TournamentStatusCompleted    TournamentStatus = "completed"
	TournamentStatusCancelled    TournamentStatus = "cancelled"
)

// Tournament representa un torneo
type Tournament struct {
	ID                  string           `json:"id"`
	Name                string           `json:"name"`
	Description         string           `json:"description"`
	TournamentType      TournamentType   `json:"tournament_type"`
	Format              TournamentFormat `json:"format"`
	EntryFee            int              `json:"entry_fee"`
	PrizePool           int              `json:"prize_pool"`
	MinRankRequired     string           `json:"min_rank_required"`
	MaxParticipants     int              `json:"max_participants"`
	CurrentParticipants int              `json:"current_participants"`
	Status              TournamentStatus `json:"status"`
	StartTime           time.Time        `json:"start_time"`
	EndTime             time.Time        `json:"end_time"`
	RegistrationStart   time.Time        `json:"registration_start"`
	RegistrationEnd     time.Time        `json:"registration_end"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// TournamentParticipant representa un participante del torneo
type TournamentParticipant struct {
	ID              string    `json:"id"`
	TournamentID    string    `json:"tournament_id"`
	UserID          string    `json:"user_id"`
	CurrentScore    int       `json:"current_score"`
	CurrentPosition int       `json:"current_position"`
	MatchesPlayed   int       `json:"matches_played"`
	MatchesWon      int       `json:"matches_won"`
	MatchesLost     int       `json:"matches_lost"`
	IsEliminated    bool      `json:"is_eliminated"`
	JoinedAt        time.Time `json:"joined_at"`
}

// TournamentMatch representa una partida del torneo
type TournamentMatch struct {
	ID            string       `json:"id"`
	TournamentID  string       `json:"tournament_id"`
	RoundNumber   int          `json:"round_number"`
	MatchNumber   int          `json:"match_number"`
	Player1ID     string       `json:"player1_id"`
	Player2ID     string       `json:"player2_id"`
	Player1Score  int          `json:"player1_score"`
	Player2Score  int          `json:"player2_score"`
	WinnerID      *string      `json:"winner_id,omitempty"`
	Status        string       `json:"status"`
	PvPMatchID    *string      `json:"pvp_match_id,omitempty"`
	ScheduledTime sql.NullTime `json:"scheduled_time,omitempty"`
	CompletedAt   sql.NullTime `json:"completed_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

// TournamentPrize representa un premio del torneo
type TournamentPrize struct {
	ID            string    `json:"id"`
	TournamentID  string    `json:"tournament_id"`
	PositionFrom  int       `json:"position_from"`
	PositionTo    int       `json:"position_to"`
	TokenReward   int       `json:"token_reward"`
	SpecialReward *string   `json:"special_reward,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// === REQUEST/RESPONSE MODELS ===

// TournamentListResponse representa la lista de torneos
type TournamentListResponse struct {
	Tournaments []TournamentWithDetails `json:"tournaments"`
	Total       int                     `json:"total"`
}

// TournamentWithDetails torneo con detalles adicionales
type TournamentWithDetails struct {
	Tournament
	Prizes           []TournamentPrize `json:"prizes"`
	IsRegistered     bool              `json:"is_registered"`
	CanRegister      bool              `json:"can_register"`
	RegistrationOpen bool              `json:"registration_open"`
	TimeUntilStart   string            `json:"time_until_start,omitempty"`
	SpotsRemaining   int               `json:"spots_remaining"`
}

// JoinTournamentRequest representa la solicitud para unirse
type JoinTournamentRequest struct {
	TournamentID string `json:"tournament_id" binding:"required"`
}

// TournamentStandingsResponse representa las posiciones del torneo
type TournamentStandingsResponse struct {
	TournamentID   string                         `json:"tournament_id"`
	TournamentName string                         `json:"tournament_name"`
	Participants   []TournamentParticipantDetails `json:"participants"`
	MyPosition     *TournamentParticipantDetails  `json:"my_position,omitempty"`
	TotalPlayers   int                            `json:"total_players"`
}

// TournamentParticipantDetails participante con detalles
type TournamentParticipantDetails struct {
	TournamentParticipant
	Username          string  `json:"username"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	RankTier          string  `json:"rank_tier"`
	IsCurrentUser     bool    `json:"is_current_user"`
}

// TournamentBracketResponse representa el bracket del torneo
type TournamentBracketResponse struct {
	TournamentID string            `json:"tournament_id"`
	Rounds       []TournamentRound `json:"rounds"`
	CurrentRound int               `json:"current_round"`
}

// TournamentRound representa una ronda del torneo
type TournamentRound struct {
	RoundNumber int                        `json:"round_number"`
	RoundName   string                     `json:"round_name"`
	Matches     []TournamentMatchWithUsers `json:"matches"`
}

// TournamentMatchWithUsers partido con info de usuarios
type TournamentMatchWithUsers struct {
	TournamentMatch
	Player1 *UserInfo `json:"player1"`
	Player2 *UserInfo `json:"player2"`
	Winner  *UserInfo `json:"winner,omitempty"`
}
