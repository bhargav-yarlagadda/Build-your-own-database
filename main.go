package main

import (
	"fmt"

	"Build-your-own-database/database/collections"
	"Build-your-own-database/database/db"
	documents "Build-your-own-database/database/document"
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
	fmt.Println("✅ Database created")

	// Step 2: Use the database
	testDB, err := dbManager.UseDatabase(dbName)
	if err != nil {
		fmt.Println("Error using database:", err)
		return
	}
	fmt.Println("✅ Using database:", dbName)

	// Step 3: Initialize CollectionManager for the selected database
	collectionManager := collections.NewCollectionManager(testDB)

	// Step 4: Create a collection
	colName := "testCollection"
	collection, err := collectionManager.CreateCollection(colName)
	if err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}
	fmt.Println("✅ Collection created:", colName)

	// Step 5: Use the collection
	collection, err = collectionManager.UseCollection(colName)
	if err != nil {
		fmt.Println("Error using collection:", err)
		return
	}
	fmt.Println("✅ Using collection:", colName)

	// Step 6: Initialize DocumentManager
	docManager := documents.NewDocumentManager(collection)

	// Step 6.1: Create Documents
	_, err = docManager.CreateDocument("doc1", map[string]interface{}{
		"name": "Alice",
		"age":  25,
	})
	if err != nil {
		fmt.Println("Error creating doc1:", err)
	} else {
		fmt.Println("✅ Created document: doc1")
	}

	_, err = docManager.CreateDocument("doc2", map[string]interface{}{
		"name": "Bob",
		"age":  30,
	})
	if err != nil {
		fmt.Println("Error creating doc2:", err)
	} else {
		fmt.Println("✅ Created document: doc2")
	}

	// Step 6.2: Use a Document
	doc, err := docManager.UseDocument("doc1")
	if err != nil {
		fmt.Println("Error using doc1:", err)
	} else {
		fmt.Printf("✅ Loaded document 'doc1': %+v\n", doc)
	}

	// Step 6.3: Find Document(s)
	results := docManager.FindDocument("name", "Bob")
	fmt.Printf("✅ Found %d document(s) with name=Bob:\n", len(results))
	for _, d := range results {
		fmt.Printf("- %s: %+v\n", d.ID, d.Data)
	}

	// Step 6.4: Delete a Document
	err = docManager.DeleteDocument("doc2")
	if err != nil {
		fmt.Println("Error deleting doc2:", err)
	} else {
		fmt.Println("✅ Deleted document: doc2")
	}

	// Step 7: Delete the collection
	err = collectionManager.DeleteCollection(colName)
	if err != nil {
		fmt.Println("Error deleting collection:", err)
		return
	}
	fmt.Println("✅ Deleted collection:", colName)

	// Step 8: Delete the database
	err = dbManager.DeleteDatabase(dbName)
	if err != nil {
		fmt.Println("Error deleting database:", err)
		return
	}
	fmt.Println("✅ Deleted database:", dbName)

	fmt.Println("🎉 All DB, collection, and document operations completed successfully!")
}
