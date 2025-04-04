package main

import (
	"fmt"
	"Build-your-own-database/database/db"
	"Build-your-own-database/database/collections"
)

func main() {
	// Initialize DBManager
	dbManager := db.NewDBManager()

	// Step 1: Create a database
	dbName := "test101"
	_, err := dbManager.CreateDatabase(dbName)
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}

	// Step 2: Use the database
	testDB, err := dbManager.UseDatabase(dbName)
	if err != nil {
		fmt.Println("Error using database:", err)
		return
	}

	// Step 3: Initialize CollectionManager for the selected database
	collectionManager := collections.NewCollectionManager(testDB)

	// Step 4: Create a collection
	colName := "testCollection"
	_, err = collectionManager.CreateCollection(colName)
	if err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}

	// Step 5: Use the collection
	_, err = collectionManager.UseCollection(colName)
	if err != nil {
		fmt.Println("Error using collection:", err)
		return
	}

	// Step 6: Delete the collection
	err = collectionManager.DeleteCollection(colName)
	if err != nil {
		fmt.Println("Error deleting collection:", err)
		return
	}

	// Step 7: Delete the database
	err = dbManager.DeleteDatabase(dbName)
	if err != nil {
		fmt.Println("Error deleting database:", err)
		return
	}

	fmt.Println("✅ All collection operations completed successfully!")
}
