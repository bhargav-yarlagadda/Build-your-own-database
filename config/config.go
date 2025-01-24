// config.go

package config

import (
	"fmt"
	"os"
)

var (
	// BasePath is the root directory where databases will be stored
	BasePath = getEnv("BASE_PATH", "C:\\Users\\bhargav\\OneDrive\\Desktop\\DatabaseStorage")
)

// getEnv is a helper function to read environment variables with a default fallback
func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// Validate checks if the required configurations are set correctly
func Validate() error {
	// You can add any other checks or validations for configuration
	if BasePath == "" {
		return fmt.Errorf("BASE_PATH is not configured")
	}
	return nil
}
