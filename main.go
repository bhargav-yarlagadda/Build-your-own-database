package main

import (
	"Build-your-own-database/database/db"
)

func main() {
	// Initialize GoDB and DBManager

	dbManager := db.NewDBManager()

	// Create a database
	dbManager.CreateDatabase("testDB")

	// Use the database
	dbManager.UseDatabase("testDB")

	// Delete the database
	dbManager.DeleteDatabase("testDB")
}
