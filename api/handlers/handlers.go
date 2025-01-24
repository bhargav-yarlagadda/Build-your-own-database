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