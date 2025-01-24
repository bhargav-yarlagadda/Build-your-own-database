package main

import (
	"fmt"
	"log"
	"Build-your-own-database/config"        // Configuration
	"Build-your-own-database/database/db"    // Database functions
	"Build-your-own-database/database/document" // Document functions
	"Build-your-own-database/database/key-value" // Key-Value functions
	// "Build-your-own-database/database/models"   // Models for Document, etc.
	"Build-your-own-database/database/utils"    // Utility functions
)

func main() {
	// Fetch configuration for the base path (assuming config/config.go defines this)
	basePath := config.BasePath

	// Step 1: Create a new database
	dbName := "TestDB"
	if utils.CheckDatabaseExists(basePath,dbName) {
		fmt.Println("Database already exists!")
	} else {
		// Create a new database if it doesn't exist
		err := database.CreateDatabase(dbName)
		if err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
		fmt.Printf("Database '%s' created successfully!\n", dbName)
	}

	// Step 2: Add a new document to the database
	docName := "TestDoc"
	if utils.CheckDocumentExists(basePath,dbName, docName) {
		fmt.Println("Document already exists!")
	} else {
		// Create a new document if it doesn't exist
		docContent := make(map[string]interface{})
		docContent["exampleKey"] = "exampleValue"

		// Add the document to the database
		_, err := document.CreateDocument(dbName, docName, docContent)
		if err != nil {
			log.Fatalf("Error adding document: %v", err)
		}
		fmt.Printf("Document '%s' added to database '%s'.\n", docName, dbName)
	}

	// Step 3: Set a key-value pair in the document
	key := "exampleKey"
	value := "newValue"
	err := keyvalues.SetKeyValue(dbName, docName, key, value)
	if err != nil {
		log.Fatalf("Error setting key-value: %v", err)
	}
	fmt.Printf("Key '%s' set to '%v' in document '%s'.\n", key, value, docName)

	// Step 4: Retrieve the key-value pair from the document
	retrievedValue, err := keyvalues.GetKeyValue(dbName, docName, key)
	if err != nil {
		log.Fatalf("Error getting key-value: %v", err)
	}
	fmt.Printf("Retrieved value for key '%s': %v\n", key, retrievedValue)

	// Step 5: Delete the key-value pair from the document
	err = keyvalues.DeleteKeyValue(dbName, docName, key)
	if err != nil {
		log.Fatalf("Error deleting key-value: %v", err)
	}
	fmt.Printf("Key '%s' deleted from document '%s'.\n", key, docName)
}
