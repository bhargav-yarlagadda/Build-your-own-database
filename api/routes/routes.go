package routes

import (
	"Build-your-own-database/api/handlers" // Import the handler
	"github.com/gofiber/fiber/v2"
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

func DeleteDocumentRoute(app *fiber.App) {

	// Route to delete a document
	app.Delete("/:databaseName/delete-document/:docName", handlers.DeleteDocumentHandler)
}
func ReadDocumentRoute(app *fiber.App){
	app.Get("/:databasename/:docname",handlers.ReadDocumentHandler)
}

func ReadAllDocumentsRoute(app *fiber.App){
	app.Get("/:dbName",handlers.ReadAllDocumentsHandler)
}
