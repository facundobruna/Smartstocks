package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIService struct {
	apiKey     string
	httpClient *http.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}

type GeneratedQuestion struct {
	Question      string `json:"question"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Category      string `json:"category"`
}

func (s *OpenAIService) GenerateQuizQuestions(difficulty string, count int) ([]GeneratedQuestion, error) {
	fmt.Println("🤖 Iniciando generación de preguntas con OpenAI...")
	fmt.Printf("   Dificultad: %s, Cantidad: %d\n", difficulty, count)

	prompt := s.buildPrompt(difficulty, count)

	reqBody := OpenAIRequest{
		Model: "openai/gpt-3.5-turbo", // OpenRouter requiere el prefijo "openai/"
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "Eres un experto en educación financiera. Genera preguntas de quiz en formato JSON válido.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	fmt.Println("📤 Enviando request a OpenAI...")

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("HTTP-Referer", "https://smartstocks.com") // OpenRouter requiere esto
	req.Header.Set("X-Title", "Smart Stocks Quiz")            // Opcional pero recomendado

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Mostrar solo los primeros caracteres disponibles
	preview := string(body)
	if len(preview) > 200 {
		preview = preview[:200]
	}
	fmt.Printf("📥 Respuesta de OpenAI (status %d): %s\n", resp.StatusCode, preview)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse el JSON de las preguntas generadas
	var questions []GeneratedQuestion
	content := openaiResp.Choices[0].Message.Content

	fmt.Printf("📝 Content recibido: %s\n", content[:min(len(content), 300)])

	if err := json.Unmarshal([]byte(content), &questions); err != nil {
		return nil, fmt.Errorf("error parsing generated questions: %w - Content: %s", err, content)
	}

	fmt.Printf("✅ %d preguntas parseadas correctamente\n", len(questions))
	return questions, nil
}

func (s *OpenAIService) buildPrompt(difficulty string, count int) string {
	difficultyDesc := map[string]string{
		"easy":   "básico, sobre conceptos fundamentales como ahorro, inversión, rendimiento fijo vs variable",
		"medium": "intermedio, sobre TIR, movimientos del mercado, análisis técnico básico",
		"hard":   "avanzado, sobre economía financiera, estrategias de inversión complejas, análisis fundamental",
	}

	return fmt.Sprintf(`Genera exactamente %d preguntas de opción múltiple sobre finanzas de nivel %s (%s).

IMPORTANTE: Responde ÚNICAMENTE con un array JSON válido, sin texto adicional antes o después.

El formato DEBE ser exactamente así:
[
  {
    "question": "Texto de la pregunta",
    "option_a": "Primera opción",
    "option_b": "Segunda opción",
    "option_c": "Tercera opción",
    "option_d": "Cuarta opción",
    "correct_answer": "A",
    "explanation": "Explicación de por qué esta es la respuesta correcta",
    "category": "Categoría de la pregunta (ej: Ahorro, Inversión, Mercados)"
  }
]

Las preguntas deben ser:
- Claras y en español
- Educativas para jóvenes argentinos
- Con 4 opciones de respuesta
- Solo una respuesta correcta (A, B, C o D)
- Con explicación clara

Genera las %d preguntas ahora:`, count, difficulty, difficultyDesc[difficulty], count)
}
