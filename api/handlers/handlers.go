package handlers

import (
	"Build-your-own-database/database/db"
	"Build-your-own-database/database/document"
	keyvalues "Build-your-own-database/database/key-value"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"github.com/gofiber/fiber/v2"
	"sync"
)

var mu sync.Mutex // Mutex to ensure concurrency is handled properly (if shared resource is involved)

func CreateDataBaseHandler(c *fiber.Ctx) error {
	var requestData struct{
		Dbname string `json:"dbname"`
	}

	// Parse the request body
	if err := c.BodyParser(&requestData); err != nil {
		log.Printf("Error in parsing Body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	// Use a goroutine to create the database asynchronously
	var wg sync.WaitGroup
	wg.Add(1) // Add one goroutine to the wait group

	go func() {
		defer wg.Done() // Ensure Done is called when goroutine finishes

		// Locking to ensure mutual exclusion for shared resources (e.g., logging or a shared db pool)
		mu.Lock()
		defer mu.Unlock()

		// Create the database
		err := db.CreateDatabase(requestData.Dbname)
		if err != nil {
			log.Printf("Error creating database: %v", err)
			// Handle error (e.g., log it or notify the user)
			return
		}

		log.Printf("Database %s created successfully!", requestData.Dbname)
	}()

	// Wait for the goroutine to finish before sending the response
	wg.Wait()

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Database created successfully!"})
}  

func DeleteDatabaseHandler(c *fiber.Ctx) error {
	// Get the database name from the request body
	var requestData struct {
		DbName string `json:"dbname"`
	}

	// Parse the request body into the struct
	if err := c.BodyParser(&requestData); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Initialize WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1) // We have one goroutine to wait for

	// Use a goroutine for the delete operation
	go func() {
		defer wg.Done() // Ensure Done is called when goroutine finishes

		// Locking to ensure mutual exclusion if needed (e.g., shared resources)
		mu.Lock()
		defer mu.Unlock()

		// Call the DeleteDatabase function from the db package
		_, err := db.DeleteDatabase(requestData.DbName)
		if err != nil {
			log.Printf("Error deleting database: %v", err)
			// Handle the error (e.g., log or notify)
			return
		}

		log.Printf("Database %s deleted successfully!", requestData.DbName)
	}()

	// Wait for the goroutine to finish before sending the response
	wg.Wait()

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Database deleted successfully!"})
}



func CreateDocumentHandler(c *fiber.Ctx) error {
	// Extract the database name from the URL parameter
	dbName := c.Params("databaseName")

	// Verify if the database exists using UseDatabase
	dbPath, err := db.UseDatabase(dbName)
	if err != nil {
		log.Printf("Database does not exist: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Database '%s' not found: %v", dbName, err)})
	}

	// Parse the request body to get the document name and content
	var requestData struct {
		DocName string                 `json:"docName"`
		Content map[string]interface{} `json:"content"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Construct the document path based on the database path
	docPath := filepath.Join(dbPath, requestData.DocName)

	// Use a WaitGroup to wait for the goroutine to complete
	var wg sync.WaitGroup
	wg.Add(1) // We have one goroutine to wait for

	// Mutex for ensuring concurrency control (if required)
	mu.Lock()
	defer mu.Unlock()

	var docUUID string

	// Run the document creation in a goroutine
	go func() {
		defer wg.Done() // Ensure Done is called when goroutine finishes

		// Call the CreateDocument function to create the document
		docUUID, err = document.CreateDocument(dbName, requestData.DocName, requestData.Content)
		if err != nil {
			log.Printf("Error creating document: %v", err)
			return
		}
	}()

	// Wait for the goroutine to finish before sending the response
	wg.Wait()

	// Check for any errors after waiting for the goroutine to finish
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create document: %v", err)})
	}

	// Return success response with the generated UUID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Document created successfully",
		"uuid":    docUUID,
		"docPath": docPath,
	})
}




// DeleteDocumentHandler handles the delete document route

func DeleteDocumentHandler(c *fiber.Ctx) error {
	// Get the database and document names from the URL params
	dbName := c.Params("databaseName")
	docName := c.Params("docName")

	// Use the database to check if it exists
	_, err := db.UseDatabase(dbName)
	if err != nil {
		log.Printf("Error verifying database: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Database not found"})
	}

	// Use a WaitGroup to wait for the goroutine to complete
	var wg sync.WaitGroup
	wg.Add(1) // We have one goroutine to wait for

	// Mutex for ensuring concurrency control (if required)
	mu.Lock()
	defer mu.Unlock()

	// Error handling variable
	var deletionError error

	// Run the document deletion in a goroutine
	go func() {
		defer wg.Done() // Ensure Done is called when goroutine finishes

		// Delete the document
		deletionError = document.DeleteDocument(dbName, docName)
		if deletionError != nil {
			log.Printf("Error deleting document: %v", deletionError)
		}
	}()

	// Wait for the goroutine to finish before sending the response
	wg.Wait()

	// Check if there was any error during document deletion
	if deletionError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete document"})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Document '%s' deleted successfully", docName)})
}
func ReadDocumentHandler(c *fiber.Ctx) error {
	dbName := c.Params("databasename")
	docName := c.Params("docname")

	// Check if the database exists
	_, err := db.UseDatabase(dbName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Database Not Found"})
	}

	// Read the document
	data, err := document.ReadDocument(dbName, docName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document Not Found"})
	}

	// Return the data as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Document retrieved successfully",
		"data":    data,
	})
}
func ReadAllDocumentsHandler(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	_, err := db.UseDatabase(dbName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Database Not Found."})
	}
	data, err := document.ReadAllDocuments(dbName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Error in Fetching Documents."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Document retrieved successfully",
		"data":    data,
	})
}


func UpdateDocumentHandler(c *fiber.Ctx) error {
	// Retrieve parameters from the request
	dbName := c.Params("dbName")
	docName := c.Params("docName")

	// Check if the database exists
	_, err := db.UseDatabase(dbName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Database Not Found."})
	}

	// Parse the request body into a map
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Use WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Mutex to ensure thread safety if the documents are shared
	mu.Lock()
	defer mu.Unlock()

	// Variable to capture any errors during the update
	var updateError error

	// Start a goroutine to handle the document update concurrently
	go func() {
		defer wg.Done() // Ensure Done is called after the goroutine finishes

		// Update the document
		updateError = document.UpdateDocument(dbName, docName, updates)
		if updateError != nil {
			log.Printf("Error updating document: %v", updateError)
		}
	}()

	// Wait for the goroutine to finish
	wg.Wait()

	// If there was an error during the update, return it
	if updateError != nil {
		if updateError.Error() == fmt.Sprintf("document '%s' does not exist", docName) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document", "details": updateError.Error()})
	}

	// Respond with success if update is successful
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Document '%s' updated successfully", docName)})
}

func DeletePairFromDocumentHandler(c *fiber.Ctx) error {
	// Extract parameters from the URL
	dbName := c.Params("dbName")
	docName := c.Params("docName")
	keyToDelete := c.Params("key")

	// Check if the key is provided
	if keyToDelete == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key to delete must be provided"})
	}

	// Call DeleteKeyValue to remove the key-value pair from the document
	err := keyvalues.DeleteKeyValue(dbName, docName, keyToDelete)
	if err != nil {
		// Check specific error cases and return appropriate responses
		if isDocumentNotFoundError(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Document '%s' not found", docName)})
		}

		if isKeyNotFoundError(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Key '%s' not found in document '%s'", keyToDelete, docName)})
		}

		// General error case for any other failures
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete key", "details": err.Error()})
	}

	// Return a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Key '%s' deleted successfully from document '%s'", keyToDelete, docName),
	})
}

// Helper function to check if the error is related to document not found
func isDocumentNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "does not exist in database")
}

// Helper function to check if the error is related to key not found
func isKeyNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "key")
}
func UpdatePairInDocumentHandler(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	docName := c.Params("docName")

	// Parse the request body to get the key and value
	var requestBody struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	// Validate key and value
	if requestBody.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key must be provided"})
	}
	if requestBody.Value == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Value must be provided"})
	}

	// Call the UpdateKeyValue function
	err := keyvalues.SetKeyValue(dbName, docName, requestBody.Key, requestBody.Value)
	if err != nil {
		// Handle specific error cases
		return handleUpdateError(err, docName, requestBody.Key)
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Key '%s' updated successfully in document '%s'", requestBody.Key, docName),
	})
}

// handleUpdateError abstracts error handling for updates, improving readability
func handleUpdateError(err error, docName, key string) error {
	// Check if document is not found
	if strings.Contains(err.Error(), "does not exist in database") {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Document '%s' not found", docName))
	}
	
	// Handle failed update cases
	if strings.Contains(err.Error(), "failed to update") {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update key '%s' in document '%s'", key, docName))
	}
	
	// General error case
	return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Unexpected error updating key '%s' in document '%s': %s", key, docName, err.Error()))
}