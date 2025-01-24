package db

import (
	"fmt"
	"os"
	"path/filepath"
	"Build-your-own-database/config"
)

// Base path to store all databases
var basePath = config.BasePath

// CreateDatabase creates a new folder for the database
func CreateDatabase(name string) error {
	// Construct the database path
	dbPath := filepath.Join(basePath, name)

	// Check if the folder already exists
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return fmt.Errorf("database '%s' already exists", name)
	}

	// Create the folder
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create database '%s': %v", name, err)
	}

	fmt.Println("Database created:", dbPath)
	return nil
}

// UseDatabase sets the active database to use
func UseDatabase(name string) (string, error) {
	// Construct the full path to the database
	dbPath := filepath.Join(basePath, name)

	// Check if the folder exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("database '%s' does not exist", name)
	}

	fmt.Println("Using database:", dbPath)
	return dbPath, nil
}

// DeleteDatabase removes a database folder and its contents
func DeleteDatabase(name string) (string, error) {
	// Construct the full path to the database
	dbPath := filepath.Join(basePath, name)

	// Check if the database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("database '%s' does not exist", name)
	}

	// Attempt to delete the database folder
	err := os.RemoveAll(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to delete database '%s': %v", name, err)
	}

	// Return success message
	return fmt.Sprintf("Database '%s' deleted successfully", name), nil
}
