package repository

import (
	"database/sql"
	"fmt"

	"github.com/smartstocks/backend/internal/models"
)

type TournamentsRepository struct {
	db *sql.DB
}

func NewTournamentsRepository(db *sql.DB) *TournamentsRepository {
	return &TournamentsRepository{db: db}
}

// GetActiveTournaments obtiene torneos activos o próximos
func (r *TournamentsRepository) GetActiveTournaments() ([]models.Tournament, error) {
	query := `
		SELECT id, name, description, tournament_type, format,
			   entry_fee, prize_pool, min_rank_required, max_participants,
			   current_participants, status, start_time, end_time,
			   registration_start, registration_end, created_at, updated_at
		FROM tournaments
		WHERE status IN ('upcoming', 'registration', 'in_progress')
		ORDER BY start_time ASC
	`

	return r.getTournaments(query)
}

// GetTournamentByID obtiene un torneo por ID
func (r *TournamentsRepository) GetTournamentByID(tournamentID string) (*models.Tournament, error) {
	tournament := &models.Tournament{}

	query := `
		SELECT id, name, description, tournament_type, format,
			   entry_fee, prize_pool, min_rank_required, max_participants,
			   current_participants, status, start_time, end_time,
			   registration_start, registration_end, created_at, updated_at
		FROM tournaments
		WHERE id = ?
	`

	err := r.db.QueryRow(query, tournamentID).Scan(
		&tournament.ID,
		&tournament.Name,
		&tournament.Description,
		&tournament.TournamentType,
		&tournament.Format,
		&tournament.EntryFee,
		&tournament.PrizePool,
		&tournament.MinRankRequired,
		&tournament.MaxParticipants,
		&tournament.CurrentParticipants,
		&tournament.Status,
		&tournament.StartTime,
		&tournament.EndTime,
		&tournament.RegistrationStart,
		&tournament.RegistrationEnd,
		&tournament.CreatedAt,
		&tournament.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tournament not found")
	}

	return tournament, err
}

// GetTournamentPrizes obtiene los premios de un torneo
func (r *TournamentsRepository) GetTournamentPrizes(tournamentID string) ([]models.TournamentPrize, error) {
	query := `
		SELECT id, tournament_id, position_from, position_to,
			   token_reward, special_reward, created_at
		FROM tournament_prizes
		WHERE tournament_id = ?
		ORDER BY position_from ASC
	`

	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []models.TournamentPrize
	for rows.Next() {
		var prize models.TournamentPrize
		var specialReward sql.NullString

		err := rows.Scan(
			&prize.ID,
			&prize.TournamentID,
			&prize.PositionFrom,
			&prize.PositionTo,
			&prize.TokenReward,
			&specialReward,
			&prize.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if specialReward.Valid {
			prize.SpecialReward = &specialReward.String
		}

		prizes = append(prizes, prize)
	}

	return prizes, nil
}

// IsUserRegistered verifica si el usuario está registrado en el torneo
func (r *TournamentsRepository) IsUserRegistered(tournamentID, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM tournament_participants WHERE tournament_id = ? AND user_id = ?`
	err := r.db.QueryRow(query, tournamentID, userID).Scan(&count)
	return count > 0, err
}

// JoinTournament inscribe a un usuario en un torneo
func (r *TournamentsRepository) JoinTournament(tournamentID, userID string) (bool, string, error) {
	var success bool
	var errorMessage sql.NullString

	query := `CALL join_tournament(?, ?, @success, @error_message)`
	_, err := r.db.Exec(query, tournamentID, userID)
	if err != nil {
		return false, "", err
	}

	err = r.db.QueryRow(`SELECT @success, @error_message`).Scan(&success, &errorMessage)
	if err != nil {
		return false, "", err
	}

	if errorMessage.Valid {
		return success, errorMessage.String, nil
	}

	return success, "", nil
}

// GetTournamentStandings obtiene las posiciones del torneo
func (r *TournamentsRepository) GetTournamentStandings(tournamentID string) ([]models.TournamentParticipant, error) {
	query := `
		SELECT id, tournament_id, user_id, current_score, current_position,
			   matches_played, matches_won, matches_lost, is_eliminated, joined_at
		FROM tournament_participants
		WHERE tournament_id = ?
		ORDER BY current_position ASC, current_score DESC
	`

	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.TournamentParticipant
	for rows.Next() {
		var p models.TournamentParticipant
		err := rows.Scan(
			&p.ID,
			&p.TournamentID,
			&p.UserID,
			&p.CurrentScore,
			&p.CurrentPosition,
			&p.MatchesPlayed,
			&p.MatchesWon,
			&p.MatchesLost,
			&p.IsEliminated,
			&p.JoinedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// GetUserParticipation obtiene la participación del usuario en un torneo
func (r *TournamentsRepository) GetUserParticipation(tournamentID, userID string) (*models.TournamentParticipant, error) {
	participant := &models.TournamentParticipant{}

	query := `
		SELECT id, tournament_id, user_id, current_score, current_position,
			   matches_played, matches_won, matches_lost, is_eliminated, joined_at
		FROM tournament_participants
		WHERE tournament_id = ? AND user_id = ?
	`

	err := r.db.QueryRow(query, tournamentID, userID).Scan(
		&participant.ID,
		&participant.TournamentID,
		&participant.UserID,
		&participant.CurrentScore,
		&participant.CurrentPosition,
		&participant.MatchesPlayed,
		&participant.MatchesWon,
		&participant.MatchesLost,
		&participant.IsEliminated,
		&participant.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return participant, err
}

// GetTournamentMatches obtiene las partidas de un torneo
func (r *TournamentsRepository) GetTournamentMatches(tournamentID string) ([]models.TournamentMatch, error) {
	query := `
		SELECT id, tournament_id, round_number, match_number,
			   player1_id, player2_id, player1_score, player2_score,
			   winner_id, status, pvp_match_id, scheduled_time, completed_at, created_at
		FROM tournament_matches
		WHERE tournament_id = ?
		ORDER BY round_number ASC, match_number ASC
	`

	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.TournamentMatch
	for rows.Next() {
		var m models.TournamentMatch
		var winnerID, pvpMatchID sql.NullString

		err := rows.Scan(
			&m.ID,
			&m.TournamentID,
			&m.RoundNumber,
			&m.MatchNumber,
			&m.Player1ID,
			&m.Player2ID,
			&m.Player1Score,
			&m.Player2Score,
			&winnerID,
			&m.Status,
			&pvpMatchID,
			&m.ScheduledTime,
			&m.CompletedAt,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if winnerID.Valid {
			m.WinnerID = &winnerID.String
		}
		if pvpMatchID.Valid {
			m.PvPMatchID = &pvpMatchID.String
		}

		matches = append(matches, m)
	}

	return matches, nil
}

// UpdateTournamentPositions actualiza las posiciones del torneo
func (r *TournamentsRepository) UpdateTournamentPositions(tournamentID string) error {
	query := `CALL update_tournament_positions(?)`
	_, err := r.db.Exec(query, tournamentID)
	return err
}

// DistributePrizes distribuye los premios del torneo
func (r *TournamentsRepository) DistributePrizes(tournamentID string) error {
	query := `CALL distribute_tournament_prizes(?)`
	_, err := r.db.Exec(query, tournamentID)
	return err
}

// GetUserTournaments obtiene los torneos en los que participa el usuario
func (r *TournamentsRepository) GetUserTournaments(userID string) ([]models.Tournament, error) {
	query := `
		SELECT DISTINCT t.id, t.name, t.description, t.tournament_type, t.format,
			   t.entry_fee, t.prize_pool, t.min_rank_required, t.max_participants,
			   t.current_participants, t.status, t.start_time, t.end_time,
			   t.registration_start, t.registration_end, t.created_at, t.updated_at
		FROM tournaments t
		JOIN tournament_participants tp ON t.id = tp.tournament_id
		WHERE tp.user_id = ?
		ORDER BY t.start_time DESC
	`

	return r.getTournaments(query, userID)
}

// Helper: getTournaments ejecuta query y devuelve torneos
func (r *TournamentsRepository) getTournaments(query string, args ...interface{}) ([]models.Tournament, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []models.Tournament
	for rows.Next() {
		var t models.Tournament
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.TournamentType,
			&t.Format,
			&t.EntryFee,
			&t.PrizePool,
			&t.MinRankRequired,
			&t.MaxParticipants,
			&t.CurrentParticipants,
			&t.Status,
			&t.StartTime,
			&t.EndTime,
			&t.RegistrationStart,
			&t.RegistrationEnd,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}

	return tournaments, nil
}
