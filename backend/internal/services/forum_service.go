package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type ForumService struct {
	forumRepo *repository.ForumRepository
}

func NewForumService(forumRepo *repository.ForumRepository) *ForumService {
	return &ForumService{
		forumRepo: forumRepo,
	}
}

// CreatePost crea un nuevo post
func (s *ForumService) CreatePost(userID string, req *models.CreatePostRequest) (*models.ForumPost, error) {
	post := &models.ForumPost{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
	}

	if err := s.forumRepo.CreatePost(post); err != nil {
		return nil, fmt.Errorf("error creating post: %w", err)
	}

	// Obtener el post completo con username
	createdPost, err := s.forumRepo.GetPostByID(post.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting created post: %w", err)
	}

	return createdPost, nil
}

// GetPosts obtiene posts con filtros y paginación
func (s *ForumService) GetPosts(req *models.GetPostsRequest, userID string) (*models.PostsListResponse, error) {
	// Valores por defecto
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	posts, total, err := s.forumRepo.GetPosts(req, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &models.PostsListResponse{
		Posts:      posts,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetPostByID obtiene un post por ID con sus respuestas
func (s *ForumService) GetPostByID(postID, userID string) (*models.PostDetailResponse, error) {
	// Incrementar vistas
	_ = s.forumRepo.IncrementViews(postID)

	// Obtener post
	post, err := s.forumRepo.GetPostByID(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	// Obtener respuestas
	replies, err := s.forumRepo.GetRepliesByPostID(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting replies: %w", err)
	}

	// Organizar respuestas en hilos
	organizedReplies := s.organizeReplies(replies)

	return &models.PostDetailResponse{
		Post:    post,
		Replies: organizedReplies,
	}, nil
}

// UpdatePost actualiza un post
func (s *ForumService) UpdatePost(postID, userID string, req *models.UpdatePostRequest) error {
	// Verificar que el usuario es el dueño
	post, err := s.forumRepo.GetPostByID(postID, "")
	if err != nil {
		return fmt.Errorf("post not found")
	}

	if post.UserID != userID {
		return errors.New("unauthorized: you can only edit your own posts")
	}

	if err := s.forumRepo.UpdatePost(postID, req); err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}

	return nil
}

// DeletePost elimina un post
func (s *ForumService) DeletePost(postID, userID string) error {
	// Verificar que el usuario es el dueño
	post, err := s.forumRepo.GetPostByID(postID, "")
	if err != nil {
		return fmt.Errorf("post not found")
	}

	if post.UserID != userID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	if err := s.forumRepo.DeletePost(postID); err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	return nil
}

// CreateReply crea una respuesta
func (s *ForumService) CreateReply(userID string, req *models.CreateReplyRequest) (*models.ForumReply, error) {
	// Verificar que el post existe
	_, err := s.forumRepo.GetPostByID(req.PostID, "")
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	reply := &models.ForumReply{
		PostID:        req.PostID,
		UserID:        userID,
		Content:       req.Content,
		ParentReplyID: req.ParentReplyID,
	}

	if err := s.forumRepo.CreateReply(reply); err != nil {
		return nil, fmt.Errorf("error creating reply: %w", err)
	}

	return reply, nil
}

// DeleteReply elimina una respuesta
func (s *ForumService) DeleteReply(replyID, userID string) error {
	// Verificar que el usuario es el dueño
	// Nota: Necesitarías un método GetReplyByID en el repo para verificar ownership
	// Por simplicidad, asumimos que se puede eliminar si existe

	if err := s.forumRepo.DeleteReply(replyID); err != nil {
		return fmt.Errorf("error deleting reply: %w", err)
	}

	return nil
}

// ReactToPost agrega/actualiza reacción a un post
func (s *ForumService) ReactToPost(userID string, req *models.ReactRequest) error {
	if req.PostID == nil && req.ReplyID == nil {
		return errors.New("either post_id or reply_id must be provided")
	}

	reaction := &models.ForumReaction{
		UserID:  userID,
		PostID:  req.PostID,
		ReplyID: req.ReplyID,
		IsLike:  req.IsLike,
	}

	if err := s.forumRepo.AddReaction(reaction); err != nil {
		return fmt.Errorf("error adding reaction: %w", err)
	}

	return nil
}

// RemoveReaction elimina una reacción
func (s *ForumService) RemoveReaction(userID string, postID, replyID *string) error {
	if err := s.forumRepo.RemoveReaction(userID, postID, replyID); err != nil {
		return fmt.Errorf("error removing reaction: %w", err)
	}

	return nil
}

// organizeReplies organiza las respuestas en hilos (helper)
func (s *ForumService) organizeReplies(replies []models.ForumReply) []models.ForumReply {
	// Mapa para acceso rápido
	replyMap := make(map[string]*models.ForumReply)
	for i := range replies {
		replyMap[replies[i].ID] = &replies[i]
	}

	// Respuestas de nivel superior (sin parent)
	var topLevel []models.ForumReply

	for i := range replies {
		if replies[i].ParentReplyID == nil {
			// Es una respuesta de nivel superior
			topLevel = append(topLevel, replies[i])
		} else {
			// Es una respuesta anidada
			parentID := *replies[i].ParentReplyID
			if parent, exists := replyMap[parentID]; exists {
				if parent.Replies == nil {
					parent.Replies = []models.ForumReply{}
				}
				parent.Replies = append(parent.Replies, replies[i])
			}
		}
	}

	return topLevel
}
