package handlers

import (
	"Build-your-own-database/database/db"
	"Build-your-own-database/database/document"
	"fmt"
	"log"
	"path/filepath"
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