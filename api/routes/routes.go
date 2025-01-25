package routes

import (
	"github.com/gofiber/fiber/v2"
	"Build-your-own-database/api/handlers" // Import the handler
)

// CreateDatabaseRoute sets up the route for creating a database
func CreateDatabaseRoute(app *fiber.App) {
	app.Post("/create-database", handlers.CreateDataBaseHandler)
}
func DeleteDatabaseRoute(app *fiber.App) {
	app.Delete("/delete-database", handlers.DeleteDatabaseHandler)
}

func CreateDocumentRoute(app *fiber.App) {
	// Route to create a document
	app.Post("/:databaseName/create-document", handlers.CreateDocumentHandler)
}