package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type SimulatorHandler struct {
	simulatorService *services.SimulatorService
}

func NewSimulatorHandler(simulatorService *services.SimulatorService) *SimulatorHandler {
	return &SimulatorHandler{
		simulatorService: simulatorService,
	}
}

// GetScenario godoc
// @Summary Get simulator scenario
// @Description Get a simulator scenario for the specified difficulty. Only one attempt per day per difficulty is allowed.
// @Tags simulator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param difficulty path string true "Difficulty level" Enums(easy, medium, hard)
// @Success 200 {object} models.SimulatorScenarioResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 429 {object} utils.Response "Already attempted today"
// @Router /simulator/{difficulty} [get]
func (h *SimulatorHandler) GetScenario(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	difficultyStr := c.Param("difficulty")
	difficulty := models.SimulatorDifficulty(difficultyStr)

	if !difficulty.IsValid() {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid difficulty. Must be: easy, medium, or hard", nil)
		return
	}

	scenario, err := h.simulatorService.GetScenario(userID, difficulty)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "you have already attempted this difficulty today. Try again tomorrow" {
			statusCode = http.StatusTooManyRequests
		}
		utils.ErrorResponse(c, statusCode, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Scenario retrieved successfully", scenario)
}

// SubmitDecision godoc
// @Summary Submit simulator decision
// @Description Submit a decision (buy/sell/hold) for a simulator scenario and get results
// @Tags simulator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.SubmitSimulatorDecisionRequest true "Submit Decision Request"
// @Success 200 {object} models.SubmitSimulatorDecisionResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /simulator/submit [post]
func (h *SimulatorHandler) SubmitDecision(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.SubmitSimulatorDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if !req.Decision.IsValid() {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid decision. Must be: buy, sell, or hold", nil)
		return
	}

	result, err := h.simulatorService.SubmitDecision(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Decision submitted successfully", result)
}

// GetHistory godoc
// @Summary Get simulator history
// @Description Get the user's simulator attempt history with statistics
// @Tags simulator
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit number of attempts (default 20, max 100)"
// @Success 200 {object} models.SimulatorHistoryResponse
// @Failure 401 {object} utils.Response
// @Router /simulator/history [get]
func (h *SimulatorHandler) GetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			limit = 20
		}
	}

	history, err := h.simulatorService.GetHistory(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get history", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "History retrieved successfully", history)
}

// GetCooldownStatus godoc
// @Summary Get cooldown status
// @Description Check if user can attempt a simulator scenario for a given difficulty
// @Tags simulator
// @Security BearerAuth
// @Produce json
// @Param difficulty path string true "Difficulty level" Enums(easy, medium, hard)
// @Success 200 {object} models.CooldownStatusResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /simulator/cooldown/{difficulty} [get]
func (h *SimulatorHandler) GetCooldownStatus(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	difficultyStr := c.Param("difficulty")
	difficulty := models.SimulatorDifficulty(difficultyStr)

	if !difficulty.IsValid() {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid difficulty. Must be: easy, medium, or hard", nil)
		return
	}

	status, err := h.simulatorService.GetCooldownStatus(userID, difficulty)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get cooldown status", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cooldown status retrieved", status)
}

// GetStats godoc
// @Summary Get simulator statistics
// @Description Get detailed statistics about the user's simulator performance
// @Tags simulator
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.SimulatorStats
// @Failure 401 {object} utils.Response
// @Router /simulator/stats [get]
func (h *SimulatorHandler) GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Obtener el historial que ya incluye las stats
	history, err := h.simulatorService.GetHistory(userID, 1) // Solo necesitamos las stats
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get stats", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stats retrieved successfully", history.Stats)
}
