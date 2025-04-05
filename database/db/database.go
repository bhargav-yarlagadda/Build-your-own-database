package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"Build-your-own-database/config"
	"Build-your-own-database/database/models"
)

type DBManager struct {
	goDB     *models.GoDB
	basePath string
	mu       sync.RWMutex
}

func NewDBManager() *DBManager {
	manager := &DBManager{
		goDB: &models.GoDB{
			Databases: make(map[string]*models.Database),
		},
		basePath: config.BasePath,
	}
	manager.loadDatabases()
	return manager
}

func (dbm *DBManager) loadDatabases() {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	entries, err := os.ReadDir(dbm.basePath)
	if err != nil {
		fmt.Println("Error reading basePath:", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dbName := entry.Name()
			dbPath := filepath.Join(dbm.basePath, dbName)

			dbm.goDB.Mutex.Lock()
			dbm.goDB.Databases[dbName] = &models.Database{
				Name:        dbName,
				Path:        dbPath,
				Collections: make(map[string]*models.Collection),
			}
			dbm.goDB.Mutex.Unlock()

			fmt.Println("Loaded database:", dbName)
		}
	}
}

func (dbm *DBManager) CreateDatabase(name string) (*models.Database, error) {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	dbm.goDB.Mutex.RLock()
	_, exists := dbm.goDB.Databases[name]
	dbm.goDB.Mutex.RUnlock()

	if exists {
		return nil, fmt.Errorf("database '%s' already exists", name)
	}

	dbPath := filepath.Join(dbm.basePath, name)
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create database '%s': %v", name, err)
	}

	db := &models.Database{
		Name:        name,
		Path:        dbPath,
		Collections: make(map[string]*models.Collection),
	}

	dbm.goDB.Mutex.Lock()
	dbm.goDB.Databases[name] = db
	dbm.goDB.Mutex.Unlock()

	fmt.Println("Database created:", name)
	return db, nil
}

func (dbm *DBManager) UseDatabase(name string) (*models.Database, error) {
	dbm.mu.RLock()
	defer dbm.mu.RUnlock()

	dbm.goDB.Mutex.RLock()
	db, exists := dbm.goDB.Databases[name]
	dbm.goDB.Mutex.RUnlock()

	if exists {
		fmt.Println("Using database:", name)
		return db, nil
	}

	dbPath := filepath.Join(dbm.basePath, name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database '%s' does not exist", name)
	}

	db = &models.Database{
		Name:        name,
		Path:        dbPath,
		Collections: make(map[string]*models.Collection),
	}

	dbm.goDB.Mutex.Lock()
	dbm.goDB.Databases[name] = db
	dbm.goDB.Mutex.Unlock()

	fmt.Println("Using database:", name)
	return db, nil
}

func (dbm *DBManager) DeleteDatabase(name string) error {
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	dbm.goDB.Mutex.Lock()
	db, exists := dbm.goDB.Databases[name]
	dbm.goDB.Mutex.Unlock()

	if !exists {
		dbPath := filepath.Join(dbm.basePath, name)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return fmt.Errorf("database '%s' does not exist", name)
		}
		db = &models.Database{
			Name: name,
			Path: filepath.Join(dbm.basePath, name),
		}
	}

	if err := os.RemoveAll(db.Path); err != nil {
		return fmt.Errorf("failed to delete database '%s': %v", name, err)
	}

	dbm.goDB.Mutex.Lock()
	delete(dbm.goDB.Databases, name)
	dbm.goDB.Mutex.Unlock()

	fmt.Println("Database deleted:", name)
	return nil
}
