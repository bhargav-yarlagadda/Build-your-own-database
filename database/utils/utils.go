package utils

import (
	"os"
	"path/filepath"
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.NewString()
}

// CheckDatabaseExists checks if the database exists at the given path
func CheckDatabaseExists(basePath, dbName string) bool {
	dbPath := filepath.Join(basePath, dbName)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CheckDocumentExists checks if the document exists in the database
func CheckDocumentExists(basePath, dbName, docName string) bool {
	docPath := filepath.Join(basePath, dbName, docName+".json")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return false
	}
	return true
}
