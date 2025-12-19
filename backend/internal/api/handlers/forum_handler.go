package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type ForumHandler struct {
	forumService *services.ForumService
}

func NewForumHandler(forumService *services.ForumService) *ForumHandler {
	return &ForumHandler{forumService: forumService}
}

// CreatePost crea un nuevo post
func (h *ForumHandler) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	post, err := h.forumService.CreatePost(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Post created successfully", post)
}

// GetPosts obtiene lista de posts con filtros
func (h *ForumHandler) GetPosts(c *gin.Context) {
	var req models.GetPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	// Obtener userID si está autenticado (opcional)
	userID, _ := middleware.GetUserID(c)

	response, err := h.forumService.GetPosts(&req, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get posts", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Posts retrieved successfully", response)
}

// GetPostByID obtiene un post específico con sus respuestas
func (h *ForumHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("id")

	// Obtener userID si está autenticado (opcional)
	userID, _ := middleware.GetUserID(c)

	response, err := h.forumService.GetPostByID(postID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post retrieved successfully", response)
}

// UpdatePost actualiza un post
func (h *ForumHandler) UpdatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	postID := c.Param("id")

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.forumService.UpdatePost(postID, userID, &req); err != nil {
		if err.Error() == "unauthorized: you can only edit your own posts" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post updated successfully", nil)
}

// DeletePost elimina un post
func (h *ForumHandler) DeletePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	postID := c.Param("id")

	if err := h.forumService.DeletePost(postID, userID); err != nil {
		if err.Error() == "unauthorized: you can only delete your own posts" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

// CreateReply crea una respuesta
func (h *ForumHandler) CreateReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	reply, err := h.forumService.CreateReply(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reply", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Reply created successfully", reply)
}

// DeleteReply elimina una respuesta
func (h *ForumHandler) DeleteReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	replyID := c.Param("id")

	if err := h.forumService.DeleteReply(replyID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reply", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reply deleted successfully", nil)
}

// ReactToPost añade o actualiza una reacción
func (h *ForumHandler) ReactToPost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.ReactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.forumService.ReactToPost(userID, &req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add reaction", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reaction added successfully", nil)
}

// RemoveReaction elimina una reacción
func (h *ForumHandler) RemoveReaction(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	postID := c.Query("post_id")
	replyID := c.Query("reply_id")

	var postIDPtr, replyIDPtr *string
	if postID != "" {
		postIDPtr = &postID
	}
	if replyID != "" {
		replyIDPtr = &replyID
	}

	if err := h.forumService.RemoveReaction(userID, postIDPtr, replyIDPtr); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove reaction", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reaction removed successfully", nil)
}
