package collections

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"Build-your-own-database/config"
	"Build-your-own-database/database/models"
)

// CollectionManager handles operations related to collections within a database
type CollectionManager struct {
	db     *models.Database
	colMux sync.Mutex
}

// NewCollectionManager initializes a CollectionManager for a given database
func NewCollectionManager(db *models.Database) *CollectionManager {
	return &CollectionManager{db: db}
}

// CreateCollection creates a new collection inside the database and persists it
func (cm *CollectionManager) CreateCollection(name string) (*models.Collection, error) {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if the collection already exists in memory
	if _, exists := cm.db.Collections[name]; exists {
		return nil, fmt.Errorf("collection '%s' already exists", name)
	}

	// Define the collection path using config.BasePath
	colPath := filepath.Join(config.BasePath, cm.db.Name, name)

	// Ensure the collection directory is created
	if err := os.MkdirAll(colPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create collection directory '%s': %v", name, err)
	}

	// Create collection object
	collection := &models.Collection{
		Name:      name,
		Documents: make(map[string]*models.Document),
	}

	// Persist collection metadata
	if err := cm.saveCollection(collection); err != nil {
		return nil, fmt.Errorf("failed to save collection metadata: %v", err)
	}

	// Store in memory
	cm.db.Collections[name] = collection

	fmt.Println("Collection created:", name)
	return collection, nil
}

// UseCollection retrieves an existing collection, loading from disk if necessary
func (cm *CollectionManager) UseCollection(name string) (*models.Collection, error) {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if the collection is already in memory
	if collection, exists := cm.db.Collections[name]; exists {
		fmt.Println("Using collection from memory:", name)
		return collection, nil
	}

	// Load collection from disk
	collection, err := cm.loadCollection(name)
	if err != nil {
		return nil, err
	}

	// Store it in memory
	cm.db.Collections[name] = collection

	fmt.Println("Loaded collection from disk:", name)
	return collection, nil
}

// DeleteCollection removes a collection from the database and disk
func (cm *CollectionManager) DeleteCollection(name string) error {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if collection exists
	if _, exists := cm.db.Collections[name]; !exists {
		return fmt.Errorf("collection '%s' does not exist", name)
	}

	// Define collection path
	colPath := filepath.Join(config.BasePath, cm.db.Name, name)

	// Delete collection directory from disk
	if err := os.RemoveAll(colPath); err != nil {
		return fmt.Errorf("failed to delete collection '%s' from disk: %v", name, err)
	}

	// Remove from memory
	delete(cm.db.Collections, name)

	fmt.Println("Collection deleted:", name)
	return nil
}

// saveCollection writes the collection metadata to a JSON file
func (cm *CollectionManager) saveCollection(collection *models.Collection) error {
	colPath := filepath.Join(config.BasePath, cm.db.Name, collection.Name, "metadata.json")
	file, err := os.Create(colPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(collection)
}

// loadCollection reads a collection from its metadata file
func (cm *CollectionManager) loadCollection(name string) (*models.Collection, error) {
	colPath := filepath.Join(config.BasePath, cm.db.Name, name, "metadata.json")

	// Check if metadata file exists
	file, err := os.Open(colPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("collection '%s' does not exist on disk", name)
		}
		return nil, err
	}
	defer file.Close()

	// Decode JSON file
	var collection models.Collection
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&collection); err != nil {
		return nil, err
	}

	// Initialize empty map if missing
	if collection.Documents == nil {
		collection.Documents = make(map[string]*models.Document)
	}

	return &collection, nil
}
