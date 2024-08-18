package utils

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
}

func GetEnvVar(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		debug.PrintStack()
		slog.Error("Environment variable not set", "envVar", name)
	}
	return envVar
}
