package keyvalues

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"Build-your-own-database/config"
	"sync"
)

var basePath = config.BasePath

// Mutex map to handle concurrency for each document
var docMutexMap = make(map[string]*sync.Mutex)

// GetOrCreateMutex returns the mutex for a given document, creating one if it doesn't exist
func GetOrCreateMutex(dbName, docName string) *sync.Mutex {
	docKey := dbName + "/" + docName
	mu, exists := docMutexMap[docKey]
	if !exists {
		mu = &sync.Mutex{}
		docMutexMap[docKey] = mu
	}
	return mu
}

// SetKeyValue sets a key-value pair in a specified document with concurrency protection
func SetKeyValue(dbName, docName, key string, value interface{}) error {
	// Get the mutex for the document and lock it
	mu := GetOrCreateMutex(dbName, docName)
	mu.Lock()
	defer mu.Unlock()

	// Construct the full path to the document
	docPath := filepath.Join(basePath, dbName, docName+".json")

	// Check if the document exists
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return fmt.Errorf("document '%s' does not exist in database '%s'", docName, dbName)
	}

	// Open the document file for reading
	file, err := os.Open(docPath)
	if err != nil {
		return fmt.Errorf("failed to open document '%s': %v", docName, err)
	}
	defer file.Close()

	// Decode the existing document content
	content := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	if err != nil {
		return fmt.Errorf("failed to decode document '%s': %v", docName, err)
	}

	// Set the key-value pair
	content[key] = value

	// Write the updated content back to the file
	file, err = os.Create(docPath) // Overwrite the file
	if err != nil {
		return fmt.Errorf("failed to write to document '%s': %v", docName, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)
	if err != nil {
		return fmt.Errorf("failed to encode content to document '%s': %v", docName, err)
	}

	fmt.Printf("Key '%s' set to '%v' in document '%s'.\n", key, value, docName)
	return nil
}

// GetKeyValue retrieves the value of a key in a specified document with concurrency protection
func GetKeyValue(dbName, docName, key string) (interface{}, error) {
	// Get the mutex for the document and lock it
	mu := GetOrCreateMutex(dbName, docName)
	mu.Lock()
	defer mu.Unlock()

	// Construct the full path to the document
	docPath := filepath.Join(basePath, dbName, docName+".json")

	// Check if the document exists
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("document '%s' does not exist in database '%s'", docName, dbName)
	}

	// Open the document file for reading
	file, err := os.Open(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document '%s': %v", docName, err)
	}
	defer file.Close()

	// Decode the existing document content
	content := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document '%s': %v", docName, err)
	}

	// Retrieve the value for the given key
	value, exists := content[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' does not exist in document '%s'", key, docName)
	}

	return value, nil
}

// DeleteKeyValue removes a key-value pair from a specified document with concurrency protection
func DeleteKeyValue(dbName, docName, key string) error {
	// Get the mutex for the document and lock it
	mu := GetOrCreateMutex(dbName, docName)
	mu.Lock()
	defer mu.Unlock()

	// Construct the full path to the document
	docPath := filepath.Join(basePath, dbName, docName+".json")

	// Check if the document exists
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return fmt.Errorf("document '%s' does not exist in database '%s'", docName, dbName)
	}

	// Open the document file for reading
	file, err := os.Open(docPath)
	if err != nil {
		return fmt.Errorf("failed to open document '%s': %v", docName, err)
	}
	defer file.Close()

	// Decode the existing document content
	content := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	if err != nil {
		return fmt.Errorf("failed to decode document '%s': %v", docName, err)
	}

	// Check if the key exists
	if _, exists := content[key]; !exists {
		return fmt.Errorf("key '%s' does not exist in document '%s'", key, docName)
	}

	// Delete the key
	delete(content, key)

	// Write the updated content back to the file
	file, err = os.Create(docPath) // Overwrite the file
	if err != nil {
		return fmt.Errorf("failed to write to document '%s': %v", docName, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)
	if err != nil {
		return fmt.Errorf("failed to encode content to document '%s': %v", docName, err)
	}

	fmt.Printf("Key '%s' deleted from document '%s'.\n", key, docName)
	return nil
}
