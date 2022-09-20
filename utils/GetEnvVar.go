package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load the .env file %s\n", err)
	}
}

func GetEnvVar(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		debug.PrintStack()
		log.Fatal(fmt.Sprintf("$%v must be set", name))
	}
	return envVar
}
