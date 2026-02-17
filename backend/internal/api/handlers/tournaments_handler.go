package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type TournamentsHandler struct {
	tournamentsService *services.TournamentsService
}

func NewTournamentsHandler(tournamentsService *services.TournamentsService) *TournamentsHandler {
	return &TournamentsHandler{
		tournamentsService: tournamentsService,
	}
}

// GetActiveTournaments godoc
// @Summary Get active tournaments
// @Description Get all active and upcoming tournaments
// @Tags tournaments
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.TournamentListResponse
// @Router /tournaments [get]
func (h *TournamentsHandler) GetActiveTournaments(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tournaments, err := h.tournamentsService.GetActiveTournaments(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get tournaments", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tournaments retrieved", tournaments)
}

// GetTournamentDetails godoc
// @Summary Get tournament details
// @Description Get detailed information about a specific tournament
// @Tags tournaments
// @Security BearerAuth
// @Produce json
// @Param tournament_id path string true "Tournament ID"
// @Success 200 {object} models.TournamentWithDetails
// @Router /tournaments/{tournament_id} [get]
func (h *TournamentsHandler) GetTournamentDetails(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tournamentID := c.Param("tournament_id")
	if tournamentID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tournament ID is required", nil)
		return
	}

	tournament, err := h.tournamentsService.GetTournamentDetails(tournamentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Tournament not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tournament details retrieved", tournament)
}

// JoinTournament godoc
// @Summary Join tournament
// @Description Register for a tournament
// @Tags tournaments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.JoinTournamentRequest true "Join Request"
// @Success 200
// @Router /tournaments/join [post]
func (h *TournamentsHandler) JoinTournament(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.JoinTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := h.tournamentsService.JoinTournament(req.TournamentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Successfully joined tournament", nil)
}

// GetTournamentStandings godoc
// @Summary Get tournament standings
// @Description Get current standings for a tournament
// @Tags tournaments
// @Security BearerAuth
// @Produce json
// @Param tournament_id path string true "Tournament ID"
// @Success 200 {object} models.TournamentStandingsResponse
// @Router /tournaments/{tournament_id}/standings [get]
func (h *TournamentsHandler) GetTournamentStandings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tournamentID := c.Param("tournament_id")
	if tournamentID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tournament ID is required", nil)
		return
	}

	standings, err := h.tournamentsService.GetTournamentStandings(tournamentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get standings", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Standings retrieved", standings)
}

// GetTournamentBracket godoc
// @Summary Get tournament bracket
// @Description Get bracket structure and matches
// @Tags tournaments
// @Security BearerAuth
// @Produce json
// @Param tournament_id path string true "Tournament ID"
// @Success 200 {object} models.TournamentBracketResponse
// @Router /tournaments/{tournament_id}/bracket [get]
func (h *TournamentsHandler) GetTournamentBracket(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tournamentID := c.Param("tournament_id")
	if tournamentID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tournament ID is required", nil)
		return
	}

	bracket, err := h.tournamentsService.GetTournamentBracket(tournamentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get bracket", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bracket retrieved", bracket)
}

// GetMyTournaments godoc
// @Summary Get my tournaments
// @Description Get all tournaments I'm participating in
// @Tags tournaments
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.TournamentListResponse
// @Router /tournaments/my-tournaments [get]
func (h *TournamentsHandler) GetMyTournaments(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	tournaments, err := h.tournamentsService.GetMyTournaments(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get tournaments", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "My tournaments retrieved", tournaments)
}
