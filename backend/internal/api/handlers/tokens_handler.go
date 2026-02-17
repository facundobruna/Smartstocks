package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/utils"
)

type TokensHandler struct {
	tokensService *services.TokensService
}

func NewTokensHandler(tokensService *services.TokensService) *TokensHandler {
	return &TokensHandler{
		tokensService: tokensService,
	}
}

// GetMyTokens godoc
// @Summary Get my tokens
// @Description Get current balance and recent transactions
// @Tags tokens
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Transaction history limit (default 20)"
// @Success 200 {object} models.TokensResponse
// @Router /tokens/balance [get]
func (h *TokensHandler) GetMyTokens(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}

	tokens, err := h.tokensService.GetUserTokens(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get tokens", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tokens retrieved", tokens)
}

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Description Get full transaction history
// @Tags tokens
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit (default 50, max 100)"
// @Success 200 {array} models.TokenTransaction
// @Router /tokens/transactions [get]
func (h *TokensHandler) GetTransactionHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 50
	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}

	transactions, err := h.tokensService.GetTransactionHistory(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get transactions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved", transactions)
}
