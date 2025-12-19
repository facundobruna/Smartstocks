package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/models"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{db: db}
}

// CreatePost crea un nuevo post
func (r *ForumRepository) CreatePost(post *models.ForumPost) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	query := `
		INSERT INTO forum_posts (id, user_id, title, content, category)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, post.ID, post.UserID, post.Title, post.Content, post.Category)
	return err
}

// GetPosts obtiene posts con filtros y paginación
func (r *ForumRepository) GetPosts(req *models.GetPostsRequest, userID string) ([]models.ForumPost, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	// Filtro por búsqueda
	if req.Search != "" {
		whereClause += " AND (title LIKE ? OR content LIKE ?)"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Filtro por categoría
	if req.Category != "" {
		whereClause += " AND category = ?"
		args = append(args, req.Category)
	}

	// Ordenamiento
	orderBy := "ORDER BY is_pinned DESC, created_at DESC"
	switch req.SortBy {
	case "likes":
		orderBy = "ORDER BY is_pinned DESC, likes DESC, created_at DESC"
	case "replies":
		orderBy = "ORDER BY is_pinned DESC, reply_count DESC, created_at DESC"
	}

	// Contar total
	countQuery := "SELECT COUNT(*) FROM forum_posts " + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Obtener posts
	offset := (req.Page - 1) * req.Limit
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.category,
			   p.likes, p.dislikes, p.reply_count, p.views, p.is_pinned, p.is_locked,
			   p.created_at, p.updated_at
		FROM forum_posts p
		JOIN users u ON p.user_id = u.id
		` + whereClause + " " + orderBy + `
		LIMIT ? OFFSET ?
	`

	args = append(args, req.Limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.ForumPost
	for rows.Next() {
		var post models.ForumPost
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.Category,
			&post.Likes, &post.Dislikes, &post.ReplyCount, &post.Views,
			&post.IsPinned, &post.IsLocked, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Obtener reacción del usuario si está autenticado
		if userID != "" {
			reaction := r.getUserReaction(userID, &post.ID, nil)
			post.UserReaction = reaction
		}

		posts = append(posts, post)
	}

	return posts, total, nil
}

// GetPostByID obtiene un post por ID
func (r *ForumRepository) GetPostByID(postID, userID string) (*models.ForumPost, error) {
	post := &models.ForumPost{}
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.category,
			   p.likes, p.dislikes, p.reply_count, p.views, p.is_pinned, p.is_locked,
			   p.created_at, p.updated_at
		FROM forum_posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`

	err := r.db.QueryRow(query, postID).Scan(
		&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.Category,
		&post.Likes, &post.Dislikes, &post.ReplyCount, &post.Views,
		&post.IsPinned, &post.IsLocked, &post.CreatedAt, &post.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, err
	}

	// Obtener reacción del usuario
	if userID != "" {
		reaction := r.getUserReaction(userID, &post.ID, nil)
		post.UserReaction = reaction
	}

	return post, nil
}

// UpdatePost actualiza un post
func (r *ForumRepository) UpdatePost(postID string, req *models.UpdatePostRequest) error {
	updates := []string{}
	args := []interface{}{}

	if req.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Content != nil {
		updates = append(updates, "content = ?")
		args = append(args, *req.Content)
	}
	if req.Category != nil {
		updates = append(updates, "category = ?")
		args = append(args, *req.Category)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, postID)

	query := "UPDATE forum_posts SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	_, err := r.db.Exec(query, args...)
	return err
}

// DeletePost elimina un post
func (r *ForumRepository) DeletePost(postID string) error {
	_, err := r.db.Exec("DELETE FROM forum_posts WHERE id = ?", postID)
	return err
}

// CreateReply crea una respuesta
func (r *ForumRepository) CreateReply(reply *models.ForumReply) error {
	reply.ID = uuid.New().String()
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = time.Now()

	query := `
		INSERT INTO forum_replies (id, post_id, user_id, parent_reply_id, content)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, reply.ID, reply.PostID, reply.UserID, reply.ParentReplyID, reply.Content)
	return err
}

// GetRepliesByPostID obtiene respuestas de un post
func (r *ForumRepository) GetRepliesByPostID(postID, userID string) ([]models.ForumReply, error) {
	query := `
		SELECT r.id, r.post_id, r.user_id, u.username, r.parent_reply_id, r.content,
			   r.likes, r.dislikes, r.is_solution, r.created_at, r.updated_at
		FROM forum_replies r
		JOIN users u ON r.user_id = u.id
		WHERE r.post_id = ?
		ORDER BY r.created_at ASC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []models.ForumReply
	for rows.Next() {
		var reply models.ForumReply
		var parentReplyID sql.NullString

		err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.UserID, &reply.Username, &parentReplyID, &reply.Content,
			&reply.Likes, &reply.Dislikes, &reply.IsSolution, &reply.CreatedAt, &reply.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if parentReplyID.Valid {
			reply.ParentReplyID = &parentReplyID.String
		}

		// Obtener reacción del usuario
		if userID != "" {
			reaction := r.getUserReaction(userID, nil, &reply.ID)
			reply.UserReaction = reaction
		}

		replies = append(replies, reply)
	}

	return replies, nil
}

// DeleteReply elimina una respuesta
func (r *ForumRepository) DeleteReply(replyID string) error {
	_, err := r.db.Exec("DELETE FROM forum_replies WHERE id = ?", replyID)
	return err
}

// AddReaction agrega o actualiza una reacción
func (r *ForumRepository) AddReaction(reaction *models.ForumReaction) error {
	reaction.ID = uuid.New().String()
	reaction.CreatedAt = time.Now()

	// Primero eliminar reacción existente
	if reaction.PostID != nil {
		r.db.Exec("DELETE FROM forum_reactions WHERE user_id = ? AND post_id = ?", reaction.UserID, *reaction.PostID)
	}
	if reaction.ReplyID != nil {
		r.db.Exec("DELETE FROM forum_reactions WHERE user_id = ? AND reply_id = ?", reaction.UserID, *reaction.ReplyID)
	}

	query := `
		INSERT INTO forum_reactions (id, user_id, post_id, reply_id, is_like)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, reaction.ID, reaction.UserID, reaction.PostID, reaction.ReplyID, reaction.IsLike)
	return err
}

// RemoveReaction elimina una reacción
func (r *ForumRepository) RemoveReaction(userID string, postID, replyID *string) error {
	if postID != nil {
		_, err := r.db.Exec("DELETE FROM forum_reactions WHERE user_id = ? AND post_id = ?", userID, *postID)
		return err
	}
	if replyID != nil {
		_, err := r.db.Exec("DELETE FROM forum_reactions WHERE user_id = ? AND reply_id = ?", userID, *replyID)
		return err
	}
	return fmt.Errorf("either post_id or reply_id must be provided")
}

// IncrementViews incrementa las vistas de un post
func (r *ForumRepository) IncrementViews(postID string) error {
	_, err := r.db.Exec("CALL increment_post_views(?)", postID)
	return err
}

// getUserReaction obtiene la reacción del usuario (helper)
func (r *ForumRepository) getUserReaction(userID string, postID, replyID *string) *string {
	var isLike bool
	var query string
	var arg interface{}

	if postID != nil {
		query = "SELECT is_like FROM forum_reactions WHERE user_id = ? AND post_id = ?"
		arg = *postID
	} else if replyID != nil {
		query = "SELECT is_like FROM forum_reactions WHERE user_id = ? AND reply_id = ?"
		arg = *replyID
	} else {
		return nil
	}

	err := r.db.QueryRow(query, userID, arg).Scan(&isLike)
	if err != nil {
		return nil
	}

	reaction := "dislike"
	if isLike {
		reaction = "like"
	}
	return &reaction
}
