package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type QuizService struct {
	quizRepo      *repository.QuizRepository
	userRepo      *repository.UserRepository
	openAIService *OpenAIService
}

func NewQuizService(
	quizRepo *repository.QuizRepository,
	userRepo *repository.UserRepository,
	openAIService *OpenAIService,
) *QuizService {
	return &QuizService{
		quizRepo:      quizRepo,
		userRepo:      userRepo,
		openAIService: openAIService,
	}
}

// GetDailyQuiz obtiene o genera el quiz diario por dificultad
func (s *QuizService) GetDailyQuiz(userID, difficulty string) (*models.QuizResponse, error) {
	// Verificar cooldown
	canAttempt, err := s.quizRepo.CheckCooldown(userID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("error checking cooldown: %w", err)
	}

	// Buscar quiz activo
	quiz, err := s.quizRepo.GetActiveQuizByDifficulty(difficulty)
	if err != nil {
		return nil, fmt.Errorf("error getting quiz: %w", err)
	}

	// Si no existe quiz o expiró, generar uno nuevo
	if quiz == nil {
		quiz, err = s.generateNewQuiz(difficulty)
		if err != nil {
			return nil, fmt.Errorf("error generating quiz: %w", err)
		}
	}

	// Obtener preguntas
	questions, err := s.quizRepo.GetQuestionsByQuizID(quiz.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting questions: %w", err)
	}

	response := &models.QuizResponse{
		Quiz:       quiz,
		Questions:  questions,
		CanAttempt: canAttempt,
	}

	// Si no puede intentar, calcular tiempo restante
	if !canAttempt {
		nextDay := time.Now().AddDate(0, 0, 1)
		midnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		hoursRemaining := int(time.Until(midnight).Hours())

		response.Cooldown = &models.CooldownInfo{
			NextAvailableAt: midnight,
			HoursRemaining:  hoursRemaining,
		}
	}

	return response, nil
}

// generateNewQuiz genera un nuevo quiz usando OpenAI
func (s *QuizService) generateNewQuiz(difficulty string) (*models.Quiz, error) {
	// Determinar puntos según dificultad
	pointsReward := map[string]int{
		"easy":   500,
		"medium": 1000,
		"hard":   2000,
	}

	// Crear quiz
	quiz := &models.Quiz{
		Difficulty:       difficulty,
		Title:            fmt.Sprintf("Quiz Diario - %s", capitalizeFirst(difficulty)),
		Description:      fmt.Sprintf("Quiz de finanzas nivel %s generado automáticamente", difficulty),
		PointsReward:     pointsReward[difficulty],
		TotalQuestions:   10,
		TimeLimitMinutes: 30,
		IsActive:         true,
		ExpiresAt:        toNullTime(time.Now().AddDate(0, 0, 1)), // Expira en 24 horas
	}

	if err := s.quizRepo.CreateQuiz(quiz); err != nil {
		return nil, fmt.Errorf("error creating quiz: %w", err)
	}

	// Generar preguntas con OpenAI
	generatedQuestions, err := s.openAIService.GenerateQuizQuestions(difficulty, 10)
	if err != nil {
		return nil, fmt.Errorf("error generating questions with AI: %w", err)
	}

	// Guardar preguntas en la base de datos
	for _, gq := range generatedQuestions {
		question := &models.QuizQuestion{
			QuizID:        quiz.ID,
			QuestionText:  gq.Question,
			OptionA:       gq.OptionA,
			OptionB:       gq.OptionB,
			OptionC:       gq.OptionC,
			OptionD:       gq.OptionD,
			CorrectOption: gq.CorrectAnswer,
			Explanation:   gq.Explanation,
			Difficulty:    difficulty,
			Category:      gq.Category,
		}

		if err := s.quizRepo.CreateQuestion(question); err != nil {
			return nil, fmt.Errorf("error saving question: %w", err)
		}
	}

	return quiz, nil
}

// SubmitQuiz procesa las respuestas del usuario
func (s *QuizService) SubmitQuiz(userID string, req *models.SubmitQuizRequest) (*models.SubmitQuizResponse, error) {
	// Verificar que el quiz existe
	questions, err := s.quizRepo.GetQuestionsByQuizID(req.QuizID)
	if err != nil {
		return nil, fmt.Errorf("error getting questions: %w", err)
	}

	if len(questions) == 0 {
		return nil, errors.New("quiz not found")
	}

	// Verificar cooldown
	difficulty := questions[0].Difficulty
	canAttempt, err := s.quizRepo.CheckCooldown(userID, difficulty)
	if err != nil {
		return nil, fmt.Errorf("error checking cooldown: %w", err)
	}

	if !canAttempt {
		return nil, errors.New("you have already completed this quiz today")
	}

	// Validar respuestas
	correctAnswers := 0
	results := []models.QuestionResult{}

	answersMap := make(map[string]string)
	for _, answer := range req.Answers {
		answersMap[answer.QuestionID] = answer.Answer
	}

	for _, question := range questions {
		userAnswer := answersMap[question.ID]
		isCorrect := userAnswer == question.CorrectOption

		if isCorrect {
			correctAnswers++
		}

		results = append(results, models.QuestionResult{
			QuestionID:    question.ID,
			QuestionText:  question.QuestionText,
			YourAnswer:    userAnswer,
			CorrectAnswer: question.CorrectOption,
			IsCorrect:     isCorrect,
			Explanation:   question.Explanation,
		})
	}

	// Calcular score y puntos
	totalQuestions := len(questions)
	score := (correctAnswers * 100) / totalQuestions

	pointsReward := map[string]int{
		"easy":   500,
		"medium": 1000,
		"hard":   2000,
	}
	pointsEarned := pointsReward[difficulty]

	// Guardar intento
	answersJSON, _ := json.Marshal(req.Answers)
	attempt := &models.QuizAttempt{
		UserID:           userID,
		QuizID:           req.QuizID,
		Difficulty:       difficulty,
		Score:            score,
		TotalQuestions:   totalQuestions,
		CorrectAnswers:   correctAnswers,
		PointsEarned:     pointsEarned,
		TimeTakenSeconds: req.TimeTakenSeconds,
		Answers:          string(answersJSON),
	}

	if err := s.quizRepo.CreateAttempt(attempt); err != nil {
		return nil, fmt.Errorf("error saving attempt: %w", err)
	}

	// Establecer cooldown
	if err := s.quizRepo.SetCooldown(userID, difficulty); err != nil {
		return nil, fmt.Errorf("error setting cooldown: %w", err)
	}

	// Actualizar stats del usuario
	if err := s.quizRepo.UpdateUserStatsAfterQuiz(userID, pointsEarned); err != nil {
		return nil, fmt.Errorf("error updating user stats: %w", err)
	}

	// Obtener stats actualizados
	userStats, err := s.userRepo.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting updated stats: %w", err)
	}

	return &models.SubmitQuizResponse{
		AttemptID:      attempt.ID,
		Score:          score,
		TotalQuestions: totalQuestions,
		CorrectAnswers: correctAnswers,
		PointsEarned:   pointsEarned,
		NewTotalPoints: userStats.Smartpoints,
		NewRank:        userStats.RankTier,
		Results:        results,
	}, nil
}

// GetQuizHistory obtiene el historial de quizzes del usuario
func (s *QuizService) GetQuizHistory(userID string) (*models.QuizHistoryResponse, error) {
	attempts, err := s.quizRepo.GetUserAttempts(userID, 50)
	if err != nil {
		return nil, fmt.Errorf("error getting attempts: %w", err)
	}

	stats, err := s.quizRepo.GetQuizStats(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting stats: %w", err)
	}

	return &models.QuizHistoryResponse{
		Attempts: attempts,
		Stats:    *stats,
	}, nil
}

// Helper functions
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
