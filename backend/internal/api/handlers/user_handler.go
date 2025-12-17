package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
	"github.com/smartstocks/backend/pkg/utils"
)

type UserHandler struct {
	userRepo   *repository.User_Repository
	schoolRepo *repository.SchoolRepository
}

func NewUserHandler(userRepo *repository.User_Repository, schoolRepo *repository.SchoolRepository) *UserHandler {
	return &UserHandler{
		userRepo:   userRepo,
		schoolRepo: schoolRepo,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved", user)
}

// GetUserStats godoc
// @Summary Get user statistics
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserStats
// @Router /user/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	stats, err := h.userRepo.GetUserStats(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Stats not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stats retrieved", stats)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "Update Request"
// @Success 200
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.userRepo.UpdateProfile(userID, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Profile update failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", nil)
}

// GetSchools godoc
// @Summary Get all schools
// @Tags user
// @Produce json
// @Success 200 {array} models.School
// @Router /schools [get]
func (h *UserHandler) GetSchools(c *gin.Context) {
	schools, err := h.schoolRepo.GetAllSchools()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve schools", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Schools retrieved", schools)
}
