package repository

import (
	"database/sql"

	"github.com/smartstocks/backend/internal/models"
)

type CoursesRepository struct {
	db *sql.DB
}

func NewCoursesRepository(db *sql.DB) *CoursesRepository {
	return &CoursesRepository{db: db}
}

// GetAllCourses obtiene todos los cursos con progreso del usuario
func (r *CoursesRepository) GetAllCourses(userID string) ([]models.Course, error) {
	query := `
		SELECT
			c.id, c.title, c.description, c.icon, c.category, c.difficulty,
			c.duration_minutes, c.points_reward, c.is_premium, c.is_active,
			c.order_index, c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM lessons WHERE course_id = c.id AND is_active = TRUE) as total_lessons,
			(SELECT COUNT(*) FROM user_lesson_progress ulp
			 JOIN lessons l ON ulp.lesson_id = l.id
			 WHERE l.course_id = c.id AND ulp.user_id = ? AND ulp.is_completed = TRUE) as completed_lessons,
			COALESCE(ucp.is_completed, FALSE) as is_completed,
			(ucp.id IS NOT NULL) as is_started
		FROM courses c
		LEFT JOIN user_course_progress ucp ON c.id = ucp.course_id AND ucp.user_id = ?
		WHERE c.is_active = TRUE
		ORDER BY c.order_index ASC
	`

	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.Icon,
			&course.Category, &course.Difficulty, &course.DurationMinutes,
			&course.PointsReward, &course.IsPremium, &course.IsActive,
			&course.OrderIndex, &course.CreatedAt, &course.UpdatedAt,
			&course.TotalLessons, &course.CompletedLessons,
			&course.IsCompleted, &course.IsStarted,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// GetCourseByID obtiene un curso por ID con progreso del usuario
func (r *CoursesRepository) GetCourseByID(courseID, userID string) (*models.Course, error) {
	query := `
		SELECT
			c.id, c.title, c.description, c.icon, c.category, c.difficulty,
			c.duration_minutes, c.points_reward, c.is_premium, c.is_active,
			c.order_index, c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM lessons WHERE course_id = c.id AND is_active = TRUE) as total_lessons,
			(SELECT COUNT(*) FROM user_lesson_progress ulp
			 JOIN lessons l ON ulp.lesson_id = l.id
			 WHERE l.course_id = c.id AND ulp.user_id = ? AND ulp.is_completed = TRUE) as completed_lessons,
			COALESCE(ucp.is_completed, FALSE) as is_completed,
			(ucp.id IS NOT NULL) as is_started
		FROM courses c
		LEFT JOIN user_course_progress ucp ON c.id = ucp.course_id AND ucp.user_id = ?
		WHERE c.id = ?
	`

	course := &models.Course{}
	err := r.db.QueryRow(query, userID, userID, courseID).Scan(
		&course.ID, &course.Title, &course.Description, &course.Icon,
		&course.Category, &course.Difficulty, &course.DurationMinutes,
		&course.PointsReward, &course.IsPremium, &course.IsActive,
		&course.OrderIndex, &course.CreatedAt, &course.UpdatedAt,
		&course.TotalLessons, &course.CompletedLessons,
		&course.IsCompleted, &course.IsStarted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return course, nil
}

// GetLessonsByCourseID obtiene las lecciones de un curso con progreso del usuario
func (r *CoursesRepository) GetLessonsByCourseID(courseID, userID string) ([]models.Lesson, error) {
	query := `
		SELECT
			l.id, l.course_id, l.title, l.content, l.content_type, l.video_url,
			l.duration_minutes, l.order_index, l.is_active, l.created_at, l.updated_at,
			COALESCE(ulp.is_completed, FALSE) as is_completed
		FROM lessons l
		LEFT JOIN user_lesson_progress ulp ON l.id = ulp.lesson_id AND ulp.user_id = ?
		WHERE l.course_id = ? AND l.is_active = TRUE
		ORDER BY l.order_index ASC
	`

	rows, err := r.db.Query(query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(
			&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
			&lesson.ContentType, &lesson.VideoURL, &lesson.DurationMinutes,
			&lesson.OrderIndex, &lesson.IsActive, &lesson.CreatedAt, &lesson.UpdatedAt,
			&lesson.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// GetLessonByID obtiene una leccion por ID con progreso del usuario
func (r *CoursesRepository) GetLessonByID(lessonID, userID string) (*models.Lesson, error) {
	query := `
		SELECT
			l.id, l.course_id, l.title, l.content, l.content_type, l.video_url,
			l.duration_minutes, l.order_index, l.is_active, l.created_at, l.updated_at,
			COALESCE(ulp.is_completed, FALSE) as is_completed
		FROM lessons l
		LEFT JOIN user_lesson_progress ulp ON l.id = ulp.lesson_id AND ulp.user_id = ?
		WHERE l.id = ?
	`

	lesson := &models.Lesson{}
	err := r.db.QueryRow(query, userID, lessonID).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
		&lesson.ContentType, &lesson.VideoURL, &lesson.DurationMinutes,
		&lesson.OrderIndex, &lesson.IsActive, &lesson.CreatedAt, &lesson.UpdatedAt,
		&lesson.IsCompleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return lesson, nil
}

// GetQuizQuestionsByLessonID obtiene las preguntas de quiz de una leccion
func (r *CoursesRepository) GetQuizQuestionsByLessonID(lessonID string) ([]models.LessonQuizQuestion, error) {
	query := `
		SELECT id, lesson_id, question_text, option_a, option_b, option_c, option_d,
			   correct_option, explanation, order_index
		FROM lesson_quiz_questions
		WHERE lesson_id = ?
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.LessonQuizQuestion
	for rows.Next() {
		var q models.LessonQuizQuestion
		err := rows.Scan(
			&q.ID, &q.LessonID, &q.QuestionText, &q.OptionA, &q.OptionB,
			&q.OptionC, &q.OptionD, &q.CorrectOption, &q.Explanation, &q.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}

// GetAdjacentLessons obtiene las lecciones anterior y siguiente
func (r *CoursesRepository) GetAdjacentLessons(courseID string, currentOrderIndex int) (*string, *string, error) {
	var prevID, nextID *string

	// Leccion anterior
	prevQuery := `
		SELECT id FROM lessons
		WHERE course_id = ? AND order_index < ? AND is_active = TRUE
		ORDER BY order_index DESC LIMIT 1
	`
	var prev string
	err := r.db.QueryRow(prevQuery, courseID, currentOrderIndex).Scan(&prev)
	if err == nil {
		prevID = &prev
	}

	// Leccion siguiente
	nextQuery := `
		SELECT id FROM lessons
		WHERE course_id = ? AND order_index > ? AND is_active = TRUE
		ORDER BY order_index ASC LIMIT 1
	`
	var next string
	err = r.db.QueryRow(nextQuery, courseID, currentOrderIndex).Scan(&next)
	if err == nil {
		nextID = &next
	}

	return prevID, nextID, nil
}

// CompleteLesson marca una leccion como completada
func (r *CoursesRepository) CompleteLesson(userID, lessonID string, quizScore *int) error {
	var score interface{} = nil
	if quizScore != nil {
		score = *quizScore
	}

	_, err := r.db.Exec("CALL complete_lesson(?, ?, ?)", userID, lessonID, score)
	return err
}

// GetCourseProgress obtiene el progreso de un curso
func (r *CoursesRepository) GetCourseProgress(userID, courseID string) (int, int, bool, error) {
	var totalLessons, completedLessons int
	var isCompleted bool

	query := `
		SELECT
			(SELECT COUNT(*) FROM lessons WHERE course_id = ? AND is_active = TRUE),
			(SELECT COUNT(*) FROM user_lesson_progress ulp
			 JOIN lessons l ON ulp.lesson_id = l.id
			 WHERE l.course_id = ? AND ulp.user_id = ? AND ulp.is_completed = TRUE),
			COALESCE((SELECT is_completed FROM user_course_progress
			          WHERE user_id = ? AND course_id = ?), FALSE)
	`

	err := r.db.QueryRow(query, courseID, courseID, userID, userID, courseID).Scan(
		&totalLessons, &completedLessons, &isCompleted,
	)
	if err != nil {
		return 0, 0, false, err
	}

	return totalLessons, completedLessons, isCompleted, nil
}

// StartCourse inicia el progreso de un curso
func (r *CoursesRepository) StartCourse(userID, courseID string) error {
	query := `
		INSERT IGNORE INTO user_course_progress (id, user_id, course_id)
		VALUES (UUID(), ?, ?)
	`
	_, err := r.db.Exec(query, userID, courseID)
	return err
}
