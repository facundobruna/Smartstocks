package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/models"
	"github.com/smartstocks/backend/internal/services"
	ws "github.com/smartstocks/backend/internal/websocket"
	"github.com/smartstocks/backend/pkg/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PvPHandler struct {
	pvpService *services.PvPService
	wsManager  *ws.Manager
}

func NewPvPHandler(pvpService *services.PvPService, wsManager *ws.Manager) *PvPHandler {
	return &PvPHandler{
		pvpService: pvpService,
		wsManager:  wsManager,
	}
}

func (h *PvPHandler) WebSocket(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		ID:      utils.GenerateID(),
		UserID:  userID,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Manager: h.wsManager,
	}

	h.wsManager.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *PvPHandler) JoinQueue(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Verificar que el usuario esté conectado por WebSocket
	if !h.wsManager.IsUserConnected(userID) {
		utils.ErrorResponse(c, http.StatusBadRequest, "You must be connected via WebSocket first", nil)
		return
	}

	log.Printf("🎮 User %s joining queue...", userID)

	response, err := h.pvpService.JoinQueue(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to join queue", err)
		return
	}

	log.Printf("✅ User %s joined queue at position %d", userID, response.Position)

	// Iniciar búsqueda de oponente en background
	go h.findMatchForUser(userID)

	utils.SuccessResponse(c, http.StatusOK, "Joined queue successfully", response)
}

func (h *PvPHandler) LeaveQueue(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.pvpService.LeaveQueue(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to leave queue", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Left queue successfully", nil)
}

func (h *PvPHandler) SubmitDecision(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.SubmitPvPDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if !req.Decision.IsValid() {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid decision. Must be: buy, sell, or hold", nil)
		return
	}

	if req.TimeElapsed < 0 || req.TimeElapsed > 15 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid time. Must be between 0 and 15 seconds", nil)
		return
	}

	log.Printf("📝 User %s submitting decision for match %s round %d: %s",
		userID, req.MatchID, req.RoundNumber, req.Decision)

	result, err := h.pvpService.SubmitDecision(userID, &req)
	if err != nil {
		if err.Error() == "waiting for opponent decision" {
			log.Printf("⏳ Waiting for opponent in match %s round %d", req.MatchID, req.RoundNumber)
			utils.SuccessResponse(c, http.StatusAccepted, "Decision submitted. Waiting for opponent...", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("✅ Round %d completed in match %s", req.RoundNumber, req.MatchID)

	// Enviar resultado por WebSocket a ambos jugadores
	h.wsManager.BroadcastToMatch(req.MatchID, result, "")

	if result.IsMatchComplete {
		log.Printf("🏆 Match %s completed!", req.MatchID)
		go h.sendMatchResult(req.MatchID, userID)
	} else {
		go h.startNextRound(req.MatchID, req.RoundNumber+1)
	}

	utils.SuccessResponse(c, http.StatusOK, "Decision submitted successfully", result)
}

func (h *PvPHandler) GetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		fmt.Sscanf(limitParam, "%d", &limit)
	}

	history, err := h.pvpService.GetHistory(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get history", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "History retrieved successfully", history)
}

// === HELPER METHODS ===

func (h *PvPHandler) findMatchForUser(userID string) {
	log.Printf("🔍 Starting matchmaking for user %s", userID)

	ticker := time.NewTicker(1 * time.Second) // Buscar cada 1 segundo
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)
	attempts := 0

	for {
		select {
		case <-ticker.C:
			attempts++

			// Verificar que el usuario siga conectado
			if !h.wsManager.IsUserConnected(userID) {
				log.Printf("❌ User %s disconnected, stopping matchmaking", userID)
				h.pvpService.LeaveQueue(userID)
				return
			}

			log.Printf("🔍 Matchmaking attempt %d for user %s", attempts, userID)

			// Buscar oponente
			matchFound, err := h.pvpService.FindMatch(userID)
			if err != nil {
				log.Printf("❌ Error finding match for user %s: %v", userID, err)
				continue
			}

			if matchFound != nil {
				// Match encontrado!
				log.Printf("🎉 MATCH FOUND! User %s vs User %s (Match ID: %s)",
					userID, matchFound.OpponentID, matchFound.MatchID)

				// Añadir ambos jugadores al match en el manager
				if client, ok := h.wsManager.GetClient(userID); ok {
					h.wsManager.AddToMatch(client, matchFound.MatchID)
					log.Printf("✅ Added user %s to match %s", userID, matchFound.MatchID)
				}
				if client, ok := h.wsManager.GetClient(matchFound.OpponentID); ok {
					h.wsManager.AddToMatch(client, matchFound.MatchID)
					log.Printf("✅ Added opponent %s to match %s", matchFound.OpponentID, matchFound.MatchID)
				}

				// Notificar a ambos jugadores
				log.Printf("📤 Sending match_found to user %s", userID)
				h.wsManager.SendToUser(userID, matchFound)

				log.Printf("📤 Sending match_found to opponent %s", matchFound.OpponentID)
				h.wsManager.SendToUser(matchFound.OpponentID, matchFound)

				// Iniciar primera ronda después de 3 segundos
				log.Printf("⏰ Starting round 1 in 3 seconds for match %s", matchFound.MatchID)
				time.Sleep(3 * time.Second)
				h.startRound(matchFound.MatchID, 1)

				return
			} else {
				log.Printf("⏳ No opponent found yet for user %s (attempt %d)", userID, attempts)
			}

		case <-timeout:
			log.Printf("⏰ Matchmaking timeout for user %s after %d attempts", userID, attempts)
			h.pvpService.LeaveQueue(userID)

			if client, ok := h.wsManager.GetClient(userID); ok {
				client.SendError("Matchmaking timeout. No opponent found.")
			}
			return
		}
	}
}

func (h *PvPHandler) startRound(matchID string, roundNumber int) {
	log.Printf("🎮 Starting round %d for match %s", roundNumber, matchID)

	roundStart, err := h.pvpService.StartRound(matchID, roundNumber)
	if err != nil {
		log.Printf("❌ Error starting round %d for match %s: %v", roundNumber, matchID, err)
		return
	}

	log.Printf("✅ Round %d started for match %s", roundNumber, matchID)
	log.Printf("📤 Broadcasting round_start to all players in match %s", matchID)

	// Enviar a todos los jugadores de la partida
	h.wsManager.BroadcastToMatch(matchID, roundStart, "")
}

func (h *PvPHandler) startNextRound(matchID string, nextRoundNumber int) {
	log.Printf("⏰ Waiting 3 seconds before starting round %d in match %s", nextRoundNumber, matchID)
	time.Sleep(3 * time.Second)
	h.startRound(matchID, nextRoundNumber)
}

func (h *PvPHandler) sendMatchResult(matchID, userID string) {
	log.Printf("⏰ Waiting 2 seconds before sending match result for %s", matchID)
	time.Sleep(2 * time.Second)

	clients := h.wsManager.GetMatchClients(matchID)
	log.Printf("📤 Sending match results to %d players in match %s", len(clients), matchID)

	for _, client := range clients {
		result, err := h.pvpService.GetMatchResult(client.UserID, matchID)
		if err != nil {
			log.Printf("❌ Error getting match result for user %s: %v", client.UserID, err)
			continue
		}

		log.Printf("✅ Sending match result to user %s: winner=%s", client.UserID, result.Winner)
		h.wsManager.SendToUser(client.UserID, result)
	}
}
