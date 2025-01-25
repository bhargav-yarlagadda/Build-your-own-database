package handlers 

import (
	"github.com/gofiber/fiber/v2"
	"Build-your-own-database/database/db"
	"log"
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