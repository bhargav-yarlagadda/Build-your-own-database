package collections

import (
	"fmt"
	"sync"

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

// CreateCollection creates a new collection inside the database
func (cm *CollectionManager) CreateCollection(name string) (*models.Collection, error) {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if the collection already exists
	if _, exists := cm.db.Collections[name]; exists {
		return nil, fmt.Errorf("collection '%s' already exists", name)
	}

	// Create a new collection
	collection := &models.Collection{
		Name:      name,
		Documents: make(map[string]*models.Document),
	}

	// Store it in the database
	cm.db.Collections[name] = collection
	fmt.Println("Collection created:", name)
	return collection, nil
}

// UseCollection retrieves an existing collection
func (cm *CollectionManager) UseCollection(name string) (*models.Collection, error) {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if the collection exists
	collection, exists := cm.db.Collections[name]
	if !exists {
		return nil, fmt.Errorf("collection '%s' does not exist", name)
	}

	fmt.Println("Using collection:", name)
	return collection, nil
}

// DeleteCollection removes a collection from the database
func (cm *CollectionManager) DeleteCollection(name string) error {
	cm.colMux.Lock()
	defer cm.colMux.Unlock()

	// Check if the collection exists
	if _, exists := cm.db.Collections[name]; !exists {
		return fmt.Errorf("collection '%s' does not exist", name)
	}

	// Delete the collection
	delete(cm.db.Collections, name)
	fmt.Println("Collection deleted:", name)
	return nil
}
