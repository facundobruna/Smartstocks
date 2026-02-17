package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type RankingsHandler struct {
	rankingsService *services.RankingsService
}

func NewRankingsHandler(rankingsService *services.RankingsService) *RankingsHandler {
	return &RankingsHandler{
		rankingsService: rankingsService,
	}
}

// GetGlobalLeaderboard godoc
// @Summary Get global leaderboard
// @Description Get the global leaderboard with top players
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit (default 100, max 1000)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} models.LeaderboardResponse
// @Router /rankings/global [get]
func (h *RankingsHandler) GetGlobalLeaderboard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}
	if offsetParam := c.Query("offset"); offsetParam != "" {
		fmt.Sscanf(offsetParam, "%d", &offset)
	}

	leaderboard, err := h.rankingsService.GetGlobalLeaderboard(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get leaderboard", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Global leaderboard retrieved", leaderboard)
}

// GetSchoolLeaderboard godoc
// @Summary Get school leaderboard
// @Description Get the leaderboard for a specific school
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Param school_id path string true "School ID"
// @Param limit query int false "Limit (default 100, max 1000)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} models.LeaderboardResponse
// @Router /rankings/school/{school_id} [get]
func (h *RankingsHandler) GetSchoolLeaderboard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	schoolID := c.Param("school_id")
	if schoolID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "School ID is required", nil)
		return
	}

	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}
	if offsetParam := c.Query("offset"); offsetParam != "" {
		fmt.Sscanf(offsetParam, "%d", &offset)
	}

	leaderboard, err := h.rankingsService.GetSchoolLeaderboard(userID, schoolID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get school leaderboard", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "School leaderboard retrieved", leaderboard)
}

// GetMySchoolLeaderboard godoc
// @Summary Get my school leaderboard
// @Description Get the leaderboard for the current user's school
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit (default 100, max 1000)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} models.LeaderboardResponse
// @Router /rankings/my-school [get]
func (h *RankingsHandler) GetMySchoolLeaderboard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}
	if offsetParam := c.Query("offset"); offsetParam != "" {
		fmt.Sscanf(offsetParam, "%d", &offset)
	}

	leaderboard, err := h.rankingsService.GetMySchoolLeaderboard(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "My school leaderboard retrieved", leaderboard)
}

// GetMyPosition godoc
// @Summary Get my ranking position
// @Description Get the current user's position in global and school rankings
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserPositionResponse
// @Router /rankings/my-position [get]
func (h *RankingsHandler) GetMyPosition(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	position, err := h.rankingsService.GetUserPosition(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get position", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Position retrieved", position)
}

// GetPublicProfile godoc
// @Summary Get public profile
// @Description Get the public profile of a user
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.UserProfilePublic
// @Router /rankings/profile/{user_id} [get]
func (h *RankingsHandler) GetPublicProfile(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	profile, err := h.rankingsService.GetPublicProfile(targetUserID, requestingUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved", profile)
}

// GetMyAchievements godoc
// @Summary Get my achievements
// @Description Get all achievements (unlocked and locked with progress)
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.AllAchievementsResponse
// @Router /rankings/achievements [get]
func (h *RankingsHandler) GetMyAchievements(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	achievements, err := h.rankingsService.GetUserAchievements(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get achievements", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Achievements retrieved", achievements)
}

// UpdateCache godoc
// @Summary Update leaderboard cache (Admin only)
// @Description Force update of the leaderboard cache
// @Tags rankings
// @Security BearerAuth
// @Produce json
// @Success 200
// @Router /rankings/admin/update-cache [post]
func (h *RankingsHandler) UpdateCache(c *gin.Context) {
	// En producción, verificar que el usuario es admin
	// Por ahora, permitir a cualquier usuario autenticado

	_, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.rankingsService.UpdateLeaderboardCache()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update cache", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cache updated successfully", nil)
}
