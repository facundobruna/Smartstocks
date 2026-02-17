package models

import (
	"database/sql"
	"time"
)

// Course representa un curso
type Course struct {
	ID              string       `json:"id"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Icon            string       `json:"icon"`
	Category        string       `json:"category"`
	Difficulty      string       `json:"difficulty"`
	DurationMinutes int          `json:"duration_minutes"`
	PointsReward    int          `json:"points_reward"`
	IsPremium       bool         `json:"is_premium"`
	IsActive        bool         `json:"is_active"`
	OrderIndex      int          `json:"order_index"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	// Campos calculados
	TotalLessons     int  `json:"total_lessons"`
	CompletedLessons int  `json:"completed_lessons"`
	IsStarted        bool `json:"is_started"`
	IsCompleted      bool `json:"is_completed"`
}

// Lesson representa una leccion de un curso
type Lesson struct {
	ID              string         `json:"id"`
	CourseID        string         `json:"course_id"`
	Title           string         `json:"title"`
	Content         string         `json:"content"`
	ContentType     string         `json:"content_type"` // text, video, quiz
	VideoURL        sql.NullString `json:"video_url,omitempty"`
	DurationMinutes int            `json:"duration_minutes"`
	OrderIndex      int            `json:"order_index"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	// Campos calculados
	IsCompleted bool `json:"is_completed"`
}

// LessonQuizQuestion representa una pregunta de quiz en una leccion
type LessonQuizQuestion struct {
	ID            string `json:"id"`
	LessonID      string `json:"lesson_id"`
	QuestionText  string `json:"question_text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectOption string `json:"correct_option"`
	Explanation   string `json:"explanation"`
	OrderIndex    int    `json:"order_index"`
}

// UserCourseProgress representa el progreso de un usuario en un curso
type UserCourseProgress struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	CourseID    string       `json:"course_id"`
	StartedAt   time.Time    `json:"started_at"`
	CompletedAt sql.NullTime `json:"completed_at,omitempty"`
	IsCompleted bool         `json:"is_completed"`
}

// UserLessonProgress representa el progreso de un usuario en una leccion
type UserLessonProgress struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	LessonID    string        `json:"lesson_id"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt sql.NullTime  `json:"completed_at,omitempty"`
	IsCompleted bool          `json:"is_completed"`
	QuizScore   sql.NullInt64 `json:"quiz_score,omitempty"`
}

// ============================================
// Request/Response types
// ============================================

// CoursesListResponse respuesta con lista de cursos
type CoursesListResponse struct {
	Courses          []Course `json:"courses"`
	TotalLessons     int      `json:"total_lessons"`
	CompletedLessons int      `json:"completed_lessons"`
	OverallProgress  float64  `json:"overall_progress"`
}

// CourseDetailResponse respuesta con detalle de un curso
type CourseDetailResponse struct {
	Course  Course   `json:"course"`
	Lessons []Lesson `json:"lessons"`
}

// LessonDetailResponse respuesta con detalle de una leccion
type LessonDetailResponse struct {
	Lesson        Lesson               `json:"lesson"`
	QuizQuestions []LessonQuizQuestion `json:"quiz_questions,omitempty"`
	NextLessonID  *string              `json:"next_lesson_id,omitempty"`
	PrevLessonID  *string              `json:"prev_lesson_id,omitempty"`
}

// CompleteLessonRequest request para completar una leccion
type CompleteLessonRequest struct {
	QuizAnswers []QuizAnswer `json:"quiz_answers,omitempty"`
}

// QuizAnswer respuesta a una pregunta del quiz
type QuizAnswer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

// CompleteLessonResponse respuesta al completar una leccion
type CompleteLessonResponse struct {
	LessonCompleted  bool `json:"lesson_completed"`
	CourseCompleted  bool `json:"course_completed"`
	QuizScore        int  `json:"quiz_score,omitempty"`
	QuizTotal        int  `json:"quiz_total,omitempty"`
	PointsEarned     int  `json:"points_earned"`
	NewTotalPoints   int  `json:"new_total_points"`
	CompletedLessons int  `json:"completed_lessons"`
	TotalLessons     int  `json:"total_lessons"`
}
