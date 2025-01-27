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
)

func CreateDataBaseHandler(c *fiber.Ctx) error {
	var requestData struct{
		Dbname string `json:"dbname"`
	}
	if err:= c.BodyParser(&requestData); err != nil{
		log.Printf("Error in parsing Body :%v",err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid body"})
	}
	err := db.CreateDatabase(requestData.Dbname)
	if err != nil{
		log.Printf("Error creating database :%v",err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create database"})
	}
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

	// Call the DeleteDatabase function from the db package
	_,err := db.DeleteDatabase(requestData.DbName)
	if err != nil {
		log.Printf("Error deleting database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete database"})
	}

	// Return a success message
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

	// Call the CreateDocument function to create the document
	docUUID, err := document.CreateDocument(dbName, requestData.DocName, requestData.Content)
	if err != nil {
		log.Printf("Error creating document: %v", err)
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

	// Delete the document
	err = document.DeleteDocument(dbName, docName)
	if err != nil {
		log.Printf("Error deleting document: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete document"})
	}

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

	// Update the document
	err = document.UpdateDocument(dbName,docName,updates)
	if err != nil {
		if err.Error() == fmt.Sprintf("document '%s' does not exist", docName) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document", "details": err.Error()})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Document '%s' updated successfully", docName)})
}

func DeletePairFromDocumentHandler(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	docName := c.Params("docName")
	keyToDelete := c.Params("key") // Assume the key to delete is passed as a parameter

	if keyToDelete == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Key to delete must be provided"})
	}

	// Call the DeleteKeyValue function
	err := keyvalues.DeleteKeyValue(dbName, docName, keyToDelete)
	if err != nil {
		// Handle specific error cases
		if strings.Contains(err.Error(), "does not exist in database") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document Not Found"})
		}
		if strings.Contains(err.Error(), "key") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete key", "details": err.Error()})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Key '%s' deleted successfully from document '%s'", keyToDelete, docName),
	})
}
