package models

import (
	"time"
)

type ForumPost struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Category   string    `json:"category"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	ReplyCount int       `json:"reply_count"`
	Views      int       `json:"views"`
	IsPinned   bool      `json:"is_pinned"`
	IsLocked   bool      `json:"is_locked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Campos adicionales para respuestas
	UserReaction *string `json:"user_reaction,omitempty"` // "like", "dislike", null
}

type ForumReply struct {
	ID            string    `json:"id"`
	PostID        string    `json:"post_id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	ParentReplyID *string   `json:"parent_reply_id"`
	Content       string    `json:"content"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
	IsSolution    bool      `json:"is_solution"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Campos adicionales
	UserReaction *string      `json:"user_reaction,omitempty"`
	Replies      []ForumReply `json:"replies,omitempty"` // Para hilos anidados
}

type ForumReaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PostID    *string   `json:"post_id,omitempty"`
	ReplyID   *string   `json:"reply_id,omitempty"`
	IsLike    bool      `json:"is_like"`
	CreatedAt time.Time `json:"created_at"`
}

// DTOs

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required,min=5,max=255"`
	Content  string `json:"content" binding:"required,min=10"`
	Category string `json:"category" binding:"omitempty,max=100"`
}

type UpdatePostRequest struct {
	Title    *string `json:"title,omitempty" binding:"omitempty,min=5,max=255"`
	Content  *string `json:"content,omitempty" binding:"omitempty,min=10"`
	Category *string `json:"category,omitempty" binding:"omitempty,max=100"`
}

type CreateReplyRequest struct {
	PostID        string  `json:"post_id" binding:"required"`
	Content       string  `json:"content" binding:"required,min=1"`
	ParentReplyID *string `json:"parent_reply_id,omitempty"`
}

type ReactRequest struct {
	PostID  *string `json:"post_id,omitempty"`
	ReplyID *string `json:"reply_id,omitempty"`
	IsLike  bool    `json:"is_like" binding:"required"`
}

type GetPostsRequest struct {
	Search   string `form:"search"`
	Category string `form:"category"`
	SortBy   string `form:"sort_by"` // "recent", "likes", "replies"
	Page     int    `form:"page" binding:"min=1"`
	Limit    int    `form:"limit" binding:"min=1,max=100"`
}

type PostDetailResponse struct {
	Post    *ForumPost   `json:"post"`
	Replies []ForumReply `json:"replies"`
}

type PostsListResponse struct {
	Posts      []ForumPost `json:"posts"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
