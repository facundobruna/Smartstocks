package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("❌ NO se cargó la API key")
	} else {
		fmt.Printf("✅ API Key cargada: %s...%s\n", apiKey[:10], apiKey[len(apiKey)-5:])
	}
}
