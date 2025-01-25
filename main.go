package main

import (
	"github.com/gofiber/fiber/v2"
	"Build-your-own-database/api/routes" // Import routes
	"log"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Set up all the routes
	routes.CreateDatabaseRoute(app)
	routes.DeleteDatabaseRoute(app) 
	routes.CreateDocumentRoute(app)
	routes.DeleteDocumentRoute(app)
	// Start the server
	log.Fatal(app.Listen(":8080"))
}


