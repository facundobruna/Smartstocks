package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type QuizHandler struct {
	quizService *services.QuizService
}

func NewQuizHandler(quizService *services.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

// GetDailyQuiz godoc
// @Summary Get daily quiz by difficulty
// @Tags quiz
// @Security BearerAuth
// @Param difficulty path string true "Difficulty" Enums(easy, medium, hard)
// @Produce json
// @Success 200 {object} models.QuizResponse
// @Router /quiz/{difficulty} [get]
func (h *QuizHandler) GetDailyQuiz(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.GetQuizRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid difficulty", err)
		return
	}

	response, err := h.quizService.GetDailyQuiz(userID, req.Difficulty)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get quiz", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz retrieved successfully", response)
}

// SubmitQuiz godoc
// @Summary Submit quiz answers
// @Tags quiz
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.SubmitQuizRequest true "Quiz Submission"
// @Success 200 {object} models.SubmitQuizResponse
// @Router /quiz/submit [post]
func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.quizService.SubmitQuiz(userID, &req)
	if err != nil {
		if err.Error() == "you have already completed this quiz today" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to submit quiz", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz submitted successfully", response)
}

// GetQuizHistory godoc
// @Summary Get user's quiz history
// @Tags quiz
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.QuizHistoryResponse
// @Router /quiz/history [get]
func (h *QuizHandler) GetQuizHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	response, err := h.quizService.GetQuizHistory(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get quiz history", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Quiz history retrieved successfully", response)
}
