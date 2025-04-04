package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"Build-your-own-database/config"
	"Build-your-own-database/database/models"
)

// DBManager struct manages database operations
type DBManager struct {
	goDB     *models.GoDB
	basePath string
	mu       sync.Mutex
}

// NewDBManager initializes the DBManager and loads databases from disk
func NewDBManager() *DBManager {
	manager := &DBManager{
		goDB: &models.GoDB{
			Databases: make(map[string]*models.Database),
		},
		basePath: config.BasePath,
	}

	// Scan and load databases
	manager.loadDatabases()

	return manager
}

// loadDatabases scans basePath and loads all databases into memory
func (dbm *DBManager) loadDatabases() {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	// Read directories under basePath
	entries, err := os.ReadDir(dbm.basePath)
	if err != nil {
		fmt.Println("Error reading basePath:", err)
		return
	}

	// Load each folder as a database
	for _, entry := range entries {
		if entry.IsDir() {
			dbName := entry.Name()
			dbm.goDB.Databases[dbName] = &models.Database{
				Name:        dbName,
				Collections: make(map[string]*models.Collection),
			}
			fmt.Println("Loaded database:", dbName)
		}
	}
}

// CreateDatabase creates a new database and stores it in memory
func (dbm *DBManager) CreateDatabase(name string) (*models.Database, error) {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	// Check if database already exists
	if _, exists := dbm.goDB.Databases[name]; exists {
		return nil, fmt.Errorf("database '%s' already exists", name)
	}

	// Create database directory
	dbPath := filepath.Join(dbm.basePath, name)
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create database '%s': %v", name, err)
	}

	// Initialize and store in memory
	db := &models.Database{
		Name:        name,
		Collections: make(map[string]*models.Collection),
	}
	dbm.goDB.Databases[name] = db

	fmt.Println("Database created:", name)
	return db, nil
}

// UseDatabase retrieves an existing database
func (dbm *DBManager) UseDatabase(name string) (*models.Database, error) {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	// Check if database exists
	db, exists := dbm.goDB.Databases[name]
	if !exists {
		// Verify if the folder exists on disk
		dbPath := filepath.Join(dbm.basePath, name)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("database '%s' does not exist", name)
		}

		// Load it into memory
		db = &models.Database{
			Name:        name,
			Collections: make(map[string]*models.Collection),
		}
		dbm.goDB.Databases[name] = db
	}

	fmt.Println("Using database:", name)
	return db, nil
}

// DeleteDatabase removes a database from memory and disk
func (dbm *DBManager) DeleteDatabase(name string) error {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	// Check if database exists in memory
	if _, exists := dbm.goDB.Databases[name]; !exists {
		// Also check if folder exists on disk
		dbPath := filepath.Join(dbm.basePath, name)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return fmt.Errorf("database '%s' does not exist", name)
		}
	}

	// Delete from disk
	dbPath := filepath.Join(dbm.basePath, name)
	err := os.RemoveAll(dbPath)
	if err != nil {
		return fmt.Errorf("failed to delete database '%s': %v", name, err)
	}

	// Remove from memory
	delete(dbm.goDB.Databases, name)

	fmt.Println("Database deleted:", name)
	return nil
}
