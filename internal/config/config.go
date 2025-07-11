package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ResolveSecret(input []byte) []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading: %v", err)
	}

	if len(input) > 0 {
		return input
	}

	if env := os.Getenv("VALIDATOR_SECRET_KEY"); env != "" {
		return []byte(env)
	}

	panic("no secret key provided and VALIDATOR_SECRET_KEY is not set")
}
