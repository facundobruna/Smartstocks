package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/smartstocks/backend/internal/models"
)

type SimulatorAIService struct {
	apiKey  string
	apiURL  string
	model   string
	timeout time.Duration
}

// NewSimulatorAIService crea un nuevo servicio de IA para el simulador
func NewSimulatorAIService(apiKey, apiURL, model string) *SimulatorAIService {
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/chat/completions"
	}
	if model == "" {
		model = "gpt-4"
	}
	return &SimulatorAIService{
		apiKey:  apiKey,
		apiURL:  apiURL,
		model:   model,
		timeout: 30 * time.Second,
	}
}

// AIScenarioRequest representa la estructura de respuesta esperada de la IA
type AIScenarioRequest struct {
	NewsContent     string                   `json:"news_content"`
	ChartData       models.ChartData         `json:"chart_data"`
	CorrectDecision models.SimulatorDecision `json:"correct_decision"`
	Explanation     string                   `json:"explanation"`
}

// GenerateScenario genera un escenario completo usando IA
func (s *SimulatorAIService) GenerateScenario(difficulty models.SimulatorDifficulty) (*models.SimulatorScenario, error) {
	prompt := s.buildPrompt(difficulty)

	response, err := s.callAI(prompt)
	if err != nil {
		return nil, fmt.Errorf("error calling AI: %w", err)
	}

	scenario, err := s.parseAIResponse(response, difficulty)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI response: %w", err)
	}

	return scenario, nil
}

// buildPrompt construye el prompt según la dificultad
func (s *SimulatorAIService) buildPrompt(difficulty models.SimulatorDifficulty) string {
	basePrompt := `Eres un experto en educación financiera para jóvenes argentinos.

TAREA: Genera un escenario de simulación de trading donde el usuario debe decidir si comprar, vender o mantener un activo.

`

	difficultyContext := ""
	switch difficulty {
	case models.SimulatorDifficultyEasy:
		difficultyContext = `DIFICULTAD: FÁCIL
- Usa conceptos básicos: acciones de empresas conocidas, criptomonedas populares
- La noticia debe ser clara y directa
- El gráfico debe mostrar una tendencia obvia (alcista o bajista)
- Ejemplos: Apple lanza nuevo iPhone, Bitcoin sube por adopción institucional
- La decisión correcta debe ser evidente`

	case models.SimulatorDifficultyMedium:
		difficultyContext = `DIFICULTAD: MEDIA
- Usa conceptos intermedios: análisis técnico básico, soportes/resistencias
- La noticia debe tener elementos mixtos (positivos y negativos)
- El gráfico debe mostrar patrones técnicos reconocibles
- Ejemplos: Tesla reporta ganancias mixtas, Mercado Libre expande operaciones
- Requiere análisis de contexto`

	case models.SimulatorDifficultyHard:
		difficultyContext = `DIFICULTAD: DIFÍCIL
- Usa conceptos avanzados: análisis fundamental, macroeconómico
- La noticia debe ser compleja con múltiples factores
- El gráfico debe mostrar señales contradictorias
- Ejemplos: Cambios en política monetaria, fusiones empresariales complejas
- Requiere análisis profundo y contra-intuitivo`
	}

	format := `
FORMATO DE RESPUESTA (DEBE SER JSON VÁLIDO):
{
  "news_content": "Noticia financiera realista y educativa (150-250 palabras)",
  "chart_data": {
    "labels": ["Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago"],
    "prices": [100, 105, 103, 108, 112, 110, 115],
    "full_prices": [100, 105, 103, 108, 112, 110, 115, 118, 122, 125, 123, 128],
    "ticker": "AAPL",
    "asset_name": "Apple Inc."
  },
  "correct_decision": "buy|sell|hold",
  "explanation": "Explicación educativa de por qué esa es la decisión correcta (100-150 palabras)"
}

REGLAS IMPORTANTES:
1. La noticia debe ser ficticia pero realista
2. 'prices' es lo que ve el usuario (7 datos)
3. 'full_prices' incluye el futuro (12 datos totales)
4. La decisión correcta debe basarse en el análisis técnico Y la noticia
5. La explicación debe ser educativa y mencionar conceptos financieros
6. Usa empresas reales o criptomonedas conocidas
7. Adapta el lenguaje para jóvenes argentinos (vos, che, etc.)
8. Los precios deben ser coherentes (no saltos irreales)
9. Incluye contexto del mercado argentino cuando sea relevante

IMPORTANTE: Responde ÚNICAMENTE con el JSON, sin markdown, sin explicaciones adicionales.`

	return basePrompt + difficultyContext + format
}

// callAI hace la llamada a la API de OpenAI
func (s *SimulatorAIService) callAI(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": s.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Eres un experto en finanzas que genera escenarios educativos. Siempre respondes en formato JSON válido.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  1500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: s.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}

// parseAIResponse parsea la respuesta de la IA y crea un escenario
func (s *SimulatorAIService) parseAIResponse(response string, difficulty models.SimulatorDifficulty) (*models.SimulatorScenario, error) {
	// Limpiar posibles markdown (```json ... ```)
	response = cleanJSONResponse(response)

	var aiScenario AIScenarioRequest
	if err := json.Unmarshal([]byte(response), &aiScenario); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w (response: %s)", err, response[:min(200, len(response))])
	}

	// Validar datos
	if aiScenario.NewsContent == "" {
		return nil, errors.New("news content is empty")
	}
	if len(aiScenario.ChartData.Prices) == 0 {
		return nil, errors.New("chart prices are empty")
	}
	if len(aiScenario.ChartData.FullPrices) <= len(aiScenario.ChartData.Prices) {
		return nil, errors.New("full prices must be longer than visible prices")
	}
	if !aiScenario.CorrectDecision.IsValid() {
		return nil, fmt.Errorf("invalid decision: %s", aiScenario.CorrectDecision)
	}

	// Crear escenario
	now := time.Now()
	scenario := &models.SimulatorScenario{
		Difficulty:      difficulty,
		NewsContent:     aiScenario.NewsContent,
		ChartData:       aiScenario.ChartData,
		CorrectDecision: aiScenario.CorrectDecision,
		Explanation:     aiScenario.Explanation,
		CreatedAt:       now,
		ExpiresAt:       now.Add(24 * time.Hour), // Expira en 24 horas
		IsActive:        true,
	}

	return scenario, nil
}

// GenerateFallbackScenario genera un escenario de respaldo si falla la IA
func (s *SimulatorAIService) GenerateFallbackScenario(difficulty models.SimulatorDifficulty) *models.SimulatorScenario {
	scenarios := s.getFallbackScenarios(difficulty)
	scenario := scenarios[rand.Intn(len(scenarios))]

	now := time.Now()
	scenario.CreatedAt = now
	scenario.ExpiresAt = now.Add(24 * time.Hour)
	scenario.IsActive = true
	scenario.Difficulty = difficulty

	return &scenario
}

// getFallbackScenarios retorna escenarios pre-generados según dificultad
func (s *SimulatorAIService) getFallbackScenarios(difficulty models.SimulatorDifficulty) []models.SimulatorScenario {
	easy := []models.SimulatorScenario{
		{
			NewsContent: "Apple anuncia resultados récord en ventas de iPhone. La compañía superó las expectativas del mercado con un aumento del 15% en ingresos trimestrales. Los analistas destacan la fuerte demanda en mercados emergentes y el éxito de los nuevos modelos con inteligencia artificial integrada.",
			ChartData: models.ChartData{
				Labels:     []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul"},
				Prices:     []float64{150, 155, 158, 162, 165, 170, 175},
				FullPrices: []float64{150, 155, 158, 162, 165, 170, 175, 180, 185, 188, 190, 195},
				Ticker:     "AAPL",
				AssetName:  "Apple Inc.",
			},
			CorrectDecision: models.SimulatorDecisionBuy,
			Explanation:     "La tendencia alcista clara, sumada a resultados financieros excelentes, indica que el momento es bueno para comprar. La empresa muestra crecimiento sostenido y las proyecciones son positivas.",
		},
	}

	medium := []models.SimulatorScenario{
		{
			NewsContent: "MercadoLibre reporta resultados mixtos. Mientras que los ingresos crecieron 20%, los costos operativos aumentaron 25%. La expansión en Brasil muestra resultados prometedores, pero la competencia en Argentina se intensifica. Los analistas están divididos sobre las perspectivas a corto plazo.",
			ChartData: models.ChartData{
				Labels:     []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul"},
				Prices:     []float64{1200, 1250, 1230, 1280, 1260, 1290, 1270},
				FullPrices: []float64{1200, 1250, 1230, 1280, 1260, 1290, 1270, 1280, 1290, 1285, 1295, 1300},
				Ticker:     "MELI",
				AssetName:  "MercadoLibre",
			},
			CorrectDecision: models.SimulatorDecisionHold,
			Explanation:     "Los resultados mixtos y la volatilidad del gráfico sugieren esperar. Aunque hay potencial de crecimiento, los costos crecientes generan incertidumbre. Mantener permite observar cómo evoluciona la situación sin riesgo adicional.",
		},
	}

	hard := []models.SimulatorScenario{
		{
			NewsContent: "El Banco Central anuncia cambios en política monetaria. La tasa de interés sube 2 puntos, buscando controlar inflación. Esto tradicionalmente fortalece la moneda pero presiona a empresas con deuda. YPF tiene alta exposición a deuda en dólares. El sector energético argentino enfrenta regulaciones nuevas que podrían afectar márgenes.",
			ChartData: models.ChartData{
				Labels:     []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul"},
				Prices:     []float64{320, 335, 328, 340, 345, 338, 350},
				FullPrices: []float64{320, 335, 328, 340, 345, 338, 350, 340, 325, 315, 305, 295},
				Ticker:     "YPF",
				AssetName:  "YPF",
			},
			CorrectDecision: models.SimulatorDecisionSell,
			Explanation:     "Aunque el gráfico muestra tendencia alcista, los factores macroeconómicos indican problemas futuros. El aumento de tasas encarece la deuda en dólares de YPF, y las nuevas regulaciones presionarán rentabilidad. Es momento de tomar ganancias antes de la corrección.",
		},
	}

	switch difficulty {
	case models.SimulatorDifficultyEasy:
		return easy
	case models.SimulatorDifficultyMedium:
		return medium
	case models.SimulatorDifficultyHard:
		return hard
	default:
		return easy
	}
}

// Helper functions

func cleanJSONResponse(response string) string {
	// Eliminar posibles markdown code blocks
	response = string(bytes.TrimPrefix([]byte(response), []byte("```json\n")))
	response = string(bytes.TrimPrefix([]byte(response), []byte("```\n")))
	response = string(bytes.TrimSuffix([]byte(response), []byte("\n```")))
	response = string(bytes.TrimSuffix([]byte(response), []byte("```")))
	return string(bytes.TrimSpace([]byte(response)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
