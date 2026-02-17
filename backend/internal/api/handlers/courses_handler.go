package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type CoursesHandler struct {
	coursesService *services.CoursesService
}

func NewCoursesHandler(coursesService *services.CoursesService) *CoursesHandler {
	return &CoursesHandler{coursesService: coursesService}
}

// GetAllCourses godoc
// @Summary Get all courses with user progress
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.CoursesListResponse
// @Router /courses [get]
func (h *CoursesHandler) GetAllCourses(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	response, err := h.coursesService.GetAllCourses(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get courses", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Courses retrieved successfully", response)
}

// GetCourseByID godoc
// @Summary Get course details with lessons
// @Tags courses
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Produce json
// @Success 200 {object} models.CourseDetailResponse
// @Router /courses/{id} [get]
func (h *CoursesHandler) GetCourseByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	courseID := c.Param("id")

	response, err := h.coursesService.GetCourseByID(courseID, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Course not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get course", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Course retrieved successfully", response)
}

// GetLessonByID godoc
// @Summary Get lesson details with quiz questions if applicable
// @Tags courses
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Produce json
// @Success 200 {object} models.LessonDetailResponse
// @Router /courses/lessons/{id} [get]
func (h *CoursesHandler) GetLessonByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	lessonID := c.Param("id")

	response, err := h.coursesService.GetLessonByID(lessonID, userID)
	if err != nil {
		if err.Error() == "lesson not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Lesson not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get lesson", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Lesson retrieved successfully", response)
}

// CompleteLesson godoc
// @Summary Complete a lesson
// @Tags courses
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Accept json
// @Produce json
// @Param request body models.CompleteLessonRequest false "Quiz answers if lesson is a quiz"
// @Success 200 {object} models.CompleteLessonResponse
// @Router /courses/lessons/{id}/complete [post]
func (h *CoursesHandler) CompleteLesson(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	lessonID := c.Param("id")

	var req models.CompleteLessonRequest
	// Quiz answers are optional
	_ = c.ShouldBindJSON(&req)

	response, err := h.coursesService.CompleteLesson(userID, lessonID, &req)
	if err != nil {
		if err.Error() == "lesson not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Lesson not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete lesson", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Lesson completed successfully", response)
}
