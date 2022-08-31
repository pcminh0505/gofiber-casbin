package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// GetEnv func to get env value
func GetEnv(key string) string {
	env := ".env"

	if os.Getenv("GO_ENV") != "" {
		env += "." + os.Getenv("GO_ENV")
	}

	err := godotenv.Load(env)
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	return os.Getenv(key)
}
