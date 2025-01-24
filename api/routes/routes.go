package routes

import (
	"github.com/gofiber/fiber/v2"
	"Build-your-own-database/api/handlers" // Import the handler
)

// CreateDatabaseRoute sets up the route for creating a database
func CreateDatabaseRoute(app *fiber.App) {
	app.Post("/create-database", handlers.CreateDataBaseHandler)
}
