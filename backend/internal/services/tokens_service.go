package services

import (
	"fmt"

	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/repository"
)

type TokensService struct {
	tokensRepo *repository.TokensRepository
}

func NewTokensService(tokensRepo *repository.TokensRepository) *TokensService {
	return &TokensService{
		tokensRepo: tokensRepo,
	}
}

// GetUserTokens obtiene el balance y transacciones recientes del usuario
func (s *TokensService) GetUserTokens(userID string, limit int) (*models.TokensResponse, error) {
	// Obtener balance
	tokens, err := s.tokensRepo.GetUserTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user tokens: %w", err)
	}

	// Obtener transacciones recientes
	transactions, err := s.tokensRepo.GetTransactionHistory(userID, limit)
	if err != nil {
		transactions = []models.TokenTransaction{} // No fallar si no hay transacciones
	}

	response := &models.TokensResponse{
		Balance:            tokens.Balance,
		TotalEarned:        tokens.TotalEarned,
		TotalSpent:         tokens.TotalSpent,
		RecentTransactions: transactions,
	}

	return response, nil
}

// GrantTokens otorga tokens a un usuario
func (s *TokensService) GrantTokens(userID string, amount int, reason, description string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return s.tokensRepo.AddTokens(userID, amount, reason, description, nil)
}

// SpendTokens gasta tokens de un usuario
func (s *TokensService) SpendTokens(userID string, amount int, reason, description string, referenceID *string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	success, err := s.tokensRepo.SubtractTokens(userID, amount, reason, description, referenceID)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("insufficient tokens")
	}

	return nil
}

// GetTransactionHistory obtiene el historial de transacciones
func (s *TokensService) GetTransactionHistory(userID string, limit int) ([]models.TokenTransaction, error) {
	return s.tokensRepo.GetTransactionHistory(userID, limit)
}
