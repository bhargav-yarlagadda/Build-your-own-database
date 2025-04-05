package main

import (
	"fmt"

	"Build-your-own-database/database/collections"
	"Build-your-own-database/database/db"
	"Build-your-own-database/database/document"
)

func main() {
	// Initialize DB Manager
	dbManager := db.NewDBManager()

	dbName := "test_db"
	colName := "test_collection"

	// Step 1: Create Database
	_, err := dbManager.CreateDatabase(dbName)
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	fmt.Println("✅ Database created")

	// Step 2: Use Database
	selectedDB, err := dbManager.UseDatabase(dbName)
	if err != nil {
		fmt.Println("Error selecting database:", err)
		return
	}
	fmt.Println("✅ Using database:", dbName)

	// Step 3: Collection Manager
	colManager := collections.NewCollectionManager(selectedDB)

	// Step 4: Create Collection
	_, err = colManager.CreateCollection(colName)
	if err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}
	fmt.Println("✅ Collection created")

	// Step 5: Use Collection
	collection, err := colManager.UseCollection(colName)
	if err != nil {
		fmt.Println("Error using collection:", err)
		return
	}
	fmt.Println("✅ Using collection:", colName)

	// Step 6: Document Manager
	docManager := documents.NewDocumentManager(collection)

	// Step 6.1: Create a document
	_, err = docManager.CreateDocument("doc1", map[string]interface{}{
		"name": "Alice",
		"age":  24,
	})
	if err != nil {
		fmt.Println("❌ Error creating doc1:", err)
		return
	}
	fmt.Println("✅ Created document: doc1")

	// Step 6.2: Use the document
	doc, err := docManager.UseDocument("doc1")
	if err != nil {
		fmt.Println("❌ Error using doc1:", err)
		return
	}

	// Step 6.3: Add field
	if err := doc.Add("email", "alice@example.com"); err != nil {
		fmt.Println("❌ Add failed:", err)
	} else {
		fmt.Println("✅ Added email field")
	}

	// Step 6.4: Find field
	if val, ok := doc.Find("email"); ok {
		fmt.Println("✅ Found email:", val)
	} else {
		fmt.Println("❌ Email not found")
	}

	// Step 6.5: Update field
	if err := doc.Update("email", "alice@new.com"); err != nil {
		fmt.Println("❌ Update failed:", err)
	} else {
		fmt.Println("✅ Updated email")
	}

	// Step 6.6: Delete a key
	if err := doc.DeleteKey("age"); err != nil {
		fmt.Println("❌ Delete key failed:", err)
	} else {
		fmt.Println("✅ Deleted age key")
	}

	// Step 6.7: Rename document
	if err := docManager.RenameDocument("doc1", "doc1_renamed"); err != nil {
		fmt.Println("❌ Rename failed:", err)
	} else {
		fmt.Println("✅ Renamed document to doc1_renamed")
	}

	// Step 6.8: Find document by value
	results := docManager.FindDocument("name", "Alice")
	fmt.Printf("✅ Found %d doc(s) with name=Alice\n", len(results))
	for _, d := range results {
		fmt.Printf("→ %s: %+v\n", d.Name, d.Data)
	}

	// Step 6.9: Delete the document
	if err := docManager.DeleteDocument("doc1_renamed"); err != nil {
		fmt.Println("❌ Document deletion failed:", err)
	} else {
		fmt.Println("✅ Deleted document: doc1_renamed")
	}

	// Step 7: Delete Collection
	if err := colManager.DeleteCollection(colName); err != nil {
		fmt.Println("❌ Collection deletion failed:", err)
	} else {
		fmt.Println("✅ Deleted collection:", colName)
	}

	// Step 8: Delete Database
	if err := dbManager.DeleteDatabase(dbName); err != nil {
		fmt.Println("❌ Database deletion failed:", err)
	} else {
		fmt.Println("✅ Deleted database:", dbName)
	}

	fmt.Println("🎉 All operations completed successfully!")
}
