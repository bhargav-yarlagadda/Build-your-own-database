package document

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"Build-your-own-database/config"
)

var basePath = config.BasePath

// CreateDocument creates a new document with an optional UUID
func CreateDocument(dbName, docName string, content map[string]interface{}) (string, error) {
	dbPath := filepath.Join(basePath, dbName)
	docPath := filepath.Join(dbPath, docName+".json")

	// Check if the database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("database %v does not exist: %v", dbName, err)
	}

	// Check if the document already exists
	if _, err := os.Stat(docPath); !os.IsNotExist(err) {
		return "", fmt.Errorf("document '%s' already exists", docName)
	}

	// Generate or use the provided UUID
	docUUID:= uuid.New().String()

	// Add the UUID to the content
	if content == nil {
		content = make(map[string]interface{})
	}
	content["uuid"] = docUUID

	// Create and write to the document file
	file, err := os.Create(docPath)
	defer file.Close()
	if err != nil {
		return "", fmt.Errorf("failed to create document '%s': %v", docName, err)
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)
	if err != nil {
		return "", fmt.Errorf("failed to write content to document '%s': %v", docName, err)
	}

	fmt.Println("Document created with UUID:", docUUID)
	return docUUID, nil
}

// ReadDocument retrieves the document using its name
func ReadDocument(dbName, docName string) (map[string]interface{}, error) {
	docPath := filepath.Join(basePath, dbName, docName+".json")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("document '%s' does not exist", docName)
	}

	file, err := os.Open(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document '%s': %v", docName, err)
	}
	defer file.Close()

	content := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document '%s': %v", docName, err)
	}

	return content, nil
}

// RetrieveDocumentByUUID fetches a document using its UUID
func RetrieveDocumentByUUID(dbName, docUUID string) (map[string]interface{}, error) {
	dbPath := filepath.Join(basePath, dbName)

	// List all files in the database directory
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database '%s': %v", dbName, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			docPath := filepath.Join(dbPath, file.Name())

			// Open and decode the file
			fileHandle, err := os.Open(docPath)
			if err != nil {
				continue // Skip if the file cannot be opened
			}
			defer fileHandle.Close()

			content := make(map[string]interface{})
			decoder := json.NewDecoder(fileHandle)
			err = decoder.Decode(&content)
			if err != nil {
				continue // Skip if decoding fails
			}

			// Check if the UUID matches
			if content["uuid"] == docUUID {
				return content, nil
			}
		}
	}

	return nil, fmt.Errorf("document with UUID '%s' not found", docUUID)
}

// UpdateDocument updates a document by its name
func UpdateDocument(dbName, docName string, updates map[string]interface{}) error {
	content, err := ReadDocument(dbName, docName)
	if err != nil {
		return err
	}

	for key, value := range updates {
		content[key] = value
	}

	docPath := filepath.Join(basePath, dbName, docName+".json")
	file, err := os.Create(docPath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to update document '%s': %v", docName, err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)
	if err != nil {
		return fmt.Errorf("failed to write updates to document '%s': %v", docName, err)
	}

	fmt.Println("Document updated:", docPath)
	return nil
}

// DeleteDocument deletes a document by its name
func DeleteDocument(dbName, docName string) error {
	docPath := filepath.Join(basePath, dbName, docName+".json")

	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return fmt.Errorf("document '%s' does not exist", docName)
	}

	err := os.Remove(docPath)
	if err != nil {
		return fmt.Errorf("failed to delete document '%s': %v", docName, err)
	}

	fmt.Println("Document deleted:", docPath)
	return nil
}

// ReadAllDocuments retrieves all documents from a specified database
func ReadAllDocuments(dbName string) ([]map[string]interface{}, error) {
	dbPath := filepath.Join(basePath, dbName)

	// Check if the database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database '%s' does not exist", dbName)
	}

	// List all files in the database directory
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database '%s': %v", dbName, err)
	}

	// Prepare a slice to hold the documents
	var documents []map[string]interface{}

	// Iterate through all files
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Open and decode the file
		docPath := filepath.Join(dbPath, file.Name())
		fileHandle, err := os.Open(docPath)
		if err != nil {
			fmt.Printf("Skipping file '%s' due to error: %v\n", file.Name(), err)
			continue
		}
		defer fileHandle.Close()

		content := make(map[string]interface{})
		decoder := json.NewDecoder(fileHandle)
		err = decoder.Decode(&content)
		if err != nil {
			fmt.Printf("Skipping file '%s' due to decoding error: %v\n", file.Name(), err)
			continue
		}

		// Add the document content to the slice
		documents = append(documents, content)
	}

	// Return the list of documents
	return documents, nil
}
