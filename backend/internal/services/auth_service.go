package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/smartstocks/backend/internal/models"

	"github.com/google/uuid"
	"github.com/smartstocks/backend/internal/repository"
	"github.com/smartstocks/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	schoolRepo       *repository.SchoolRepository
	jwtManager       *jwt.JWTManager
	refreshTokenDays int
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	schoolRepo *repository.SchoolRepository,
	jwtManager *jwt.JWTManager,
	refreshTokenDays int,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		schoolRepo:       schoolRepo,
		jwtManager:       jwtManager,
		refreshTokenDays: refreshTokenDays,
	}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validar si el email ya existe
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Validar si el username ya existe
	existingUsername, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil && existingUsername != nil {
		return nil, errors.New("username already taken")
	}

	// Validar school_id si se proporciona
	if req.SchoolID != nil {
		school, err := s.schoolRepo.GetSchoolByID(*req.SchoolID)
		if err != nil || school == nil {
			return nil, errors.New("invalid school_id")
		}
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Generar token de verificación
	verificationToken := uuid.New().String()

	// Crear usuario
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		EmailVerified: false,
		VerificationToken: sql.NullString{
			String: verificationToken,
			Valid:  true,
		},
	}

	if req.ProfilePictureURL != nil {
		user.ProfilePictureURL = sql.NullString{String: *req.ProfilePictureURL, Valid: true}
	}

	if req.SchoolID != nil {
		user.SchoolID = sql.NullString{String: *req.SchoolID, Valid: true}
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// TODO: Enviar email de verificación con el token
	// s.emailService.SendVerificationEmail(user.Email, verificationToken)

	// Generar tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken := s.jwtManager.GenerateRefreshToken()
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: s.jwtManager.GetRefreshTokenExpiration(s.refreshTokenDays),
	}

	if err := s.refreshTokenRepo.CreateRefreshToken(refreshTokenModel); err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	// Obtener stats del usuario
	stats, _ := s.userRepo.GetUserStats(user.ID)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.userToUserInfo(user),
		Stats:        stats,
	}, nil
}

// Login autentica un usuario
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Buscar usuario
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Actualizar último login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generar tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken := s.jwtManager.GenerateRefreshToken()
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: s.jwtManager.GetRefreshTokenExpiration(s.refreshTokenDays),
	}

	if err := s.refreshTokenRepo.CreateRefreshToken(refreshTokenModel); err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	// Obtener stats
	stats, _ := s.userRepo.GetUserStats(user.ID)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.userToUserInfo(user),
		Stats:        stats,
	}, nil
}

// RefreshToken genera un nuevo access token usando un refresh token
func (s *AuthService) RefreshToken(refreshTokenStr string) (*models.LoginResponse, error) {
	// Validar refresh token
	refreshToken, err := s.refreshTokenRepo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Obtener usuario
	user, err := s.userRepo.GetUserByID(refreshToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generar nuevo access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Obtener stats
	stats, _ := s.userRepo.GetUserStats(user.ID)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         s.userToUserInfo(user),
		Stats:        stats,
	}, nil
}

// VerifyEmail verifica el email de un usuario
func (s *AuthService) VerifyEmail(token string) error {
	return s.userRepo.VerifyEmail(token)
}

// Logout invalida el refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.DeleteRefreshToken(refreshToken)
}

// Helper: convierte User a UserInfo
func (s *AuthService) userToUserInfo(user *models.User) *models.UserInfo {
	userInfo := &models.UserInfo{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.ProfilePictureURL.Valid {
		userInfo.ProfilePictureURL = &user.ProfilePictureURL.String
	}

	if user.SchoolID.Valid {
		userInfo.SchoolID = &user.SchoolID.String
	}

	return userInfo
}
