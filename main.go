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
	routes.ReadDocumentRoute(app)
	routes.DeleteDocumentRoute(app)
	routes.ReadAllDocumentsRoute(app)
	routes.UpdateDocumentRoute(app)
	routes.DeletePairFromDocumentRoute(app)
	// Start the server
	log.Fatal(app.Listen(":8080"))
}


