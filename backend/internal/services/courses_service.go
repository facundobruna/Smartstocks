package services

import (
	"errors"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type CoursesService struct {
	coursesRepo *repository.CoursesRepository
	userRepo    *repository.UserRepository
}

func NewCoursesService(coursesRepo *repository.CoursesRepository, userRepo *repository.UserRepository) *CoursesService {
	return &CoursesService{
		coursesRepo: coursesRepo,
		userRepo:    userRepo,
	}
}

// GetAllCourses obtiene todos los cursos con progreso del usuario
func (s *CoursesService) GetAllCourses(userID string) (*models.CoursesListResponse, error) {
	courses, err := s.coursesRepo.GetAllCourses(userID)
	if err != nil {
		return nil, err
	}

	// Calcular totales
	var totalLessons, completedLessons int
	for _, course := range courses {
		totalLessons += course.TotalLessons
		completedLessons += course.CompletedLessons
	}

	var overallProgress float64
	if totalLessons > 0 {
		overallProgress = float64(completedLessons) / float64(totalLessons) * 100
	}

	return &models.CoursesListResponse{
		Courses:          courses,
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
		OverallProgress:  overallProgress,
	}, nil
}

// GetCourseByID obtiene un curso con sus lecciones
func (s *CoursesService) GetCourseByID(courseID, userID string) (*models.CourseDetailResponse, error) {
	course, err := s.coursesRepo.GetCourseByID(courseID, userID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	lessons, err := s.coursesRepo.GetLessonsByCourseID(courseID, userID)
	if err != nil {
		return nil, err
	}

	return &models.CourseDetailResponse{
		Course:  *course,
		Lessons: lessons,
	}, nil
}

// GetLessonByID obtiene una leccion con sus preguntas de quiz
func (s *CoursesService) GetLessonByID(lessonID, userID string) (*models.LessonDetailResponse, error) {
	lesson, err := s.coursesRepo.GetLessonByID(lessonID, userID)
	if err != nil {
		return nil, err
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}

	response := &models.LessonDetailResponse{
		Lesson: *lesson,
	}

	// Si es un quiz, obtener preguntas
	if lesson.ContentType == "quiz" {
		questions, err := s.coursesRepo.GetQuizQuestionsByLessonID(lessonID)
		if err != nil {
			return nil, err
		}
		response.QuizQuestions = questions
	}

	// Obtener lecciones adyacentes
	prevID, nextID, err := s.coursesRepo.GetAdjacentLessons(lesson.CourseID, lesson.OrderIndex)
	if err != nil {
		return nil, err
	}
	response.PrevLessonID = prevID
	response.NextLessonID = nextID

	// Iniciar progreso del curso si no existe
	_ = s.coursesRepo.StartCourse(userID, lesson.CourseID)

	return response, nil
}

// CompleteLesson completa una leccion
func (s *CoursesService) CompleteLesson(userID, lessonID string, req *models.CompleteLessonRequest) (*models.CompleteLessonResponse, error) {
	lesson, err := s.coursesRepo.GetLessonByID(lessonID, userID)
	if err != nil {
		return nil, err
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}

	var quizScore *int
	var quizTotal int

	// Si es un quiz, validar respuestas
	if lesson.ContentType == "quiz" && req != nil && len(req.QuizAnswers) > 0 {
		questions, err := s.coursesRepo.GetQuizQuestionsByLessonID(lessonID)
		if err != nil {
			return nil, err
		}

		correctAnswers := 0
		questionMap := make(map[string]string)
		for _, q := range questions {
			questionMap[q.ID] = q.CorrectOption
		}

		for _, answer := range req.QuizAnswers {
			if correct, exists := questionMap[answer.QuestionID]; exists {
				if answer.Answer == correct {
					correctAnswers++
				}
			}
		}

		quizTotal = len(questions)
		score := correctAnswers
		quizScore = &score
	}

	// Completar la leccion
	err = s.coursesRepo.CompleteLesson(userID, lessonID, quizScore)
	if err != nil {
		return nil, err
	}

	// Obtener progreso actualizado
	totalLessons, completedLessons, courseCompleted, err := s.coursesRepo.GetCourseProgress(userID, lesson.CourseID)
	if err != nil {
		return nil, err
	}

	// Obtener puntos actualizados
	userStats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, err
	}

	// Calcular puntos ganados
	var pointsEarned int
	if courseCompleted {
		course, _ := s.coursesRepo.GetCourseByID(lesson.CourseID, userID)
		if course != nil {
			pointsEarned = course.PointsReward
		}
	}

	response := &models.CompleteLessonResponse{
		LessonCompleted:  true,
		CourseCompleted:  courseCompleted,
		PointsEarned:     pointsEarned,
		NewTotalPoints:   userStats.Smartpoints,
		CompletedLessons: completedLessons,
		TotalLessons:     totalLessons,
	}

	if quizScore != nil {
		response.QuizScore = *quizScore
		response.QuizTotal = quizTotal
	}

	return response, nil
}
