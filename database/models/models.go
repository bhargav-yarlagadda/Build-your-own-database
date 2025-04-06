package models

import "sync"

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)


// GoDB is the central database manager
type GoDB struct {
	Databases map[string]*Database // Stores all databases
	Mutex     sync.RWMutex         // Protects access to Databases
}

// Database represents a database in the system
type Database struct {
	Name        string                  `json:"name"`        // Name of the database
	Path        string                  `json:"path"`        // Path where the database is stored
	Collections map[string]*Collection  `json:"collections"` // List of collections in the database
	Mutex       sync.RWMutex            // Protects access to Collections
}

// Collection represents a collection inside a database
type Collection struct {
	Name      string                   `json:"name"`
	Path      string                   `json:"path"`
	Documents map[string]*Document     `json:"documents"`

	mu sync.Mutex `json:"-"` // Prevent mutex from being serialized
}


// Document represents an individual document inside a collection
type Document struct {
	ID   string                 `json:"id"`   // Document ID
	Name string 				`json:"name"`
	Data map[string]interface{} `json:"data"` // Key-value data
	Path string                 `json:"path"` // Path to the file on disk (optional)
}


func (d *Document) Add(key string, value interface{}) error {
	if _, exists := d.Data[key]; exists {
		return fmt.Errorf("key '%s' already exists", key)
	}
	d.Data[key] = value
	return d.save()
}

func (d *Document) Find(key string) (interface{}, bool) {
	val, ok := d.Data[key]
	return val, ok
}

func (d *Document) Update(key string, value interface{}) error {
	if _, exists := d.Data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}
	d.Data[key] = value
	return d.save()
}

func (d *Document) DeleteKey(key string) error {
	if _, exists := d.Data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}
	delete(d.Data, key)
	return d.save()
}

func (d *Document) Rename(newID string) error {
	newPath := filepath.Join(filepath.Dir(d.Path), newID+".json")
	if err := os.Rename(d.Path, newPath); err != nil {
		return err
	}
	d.ID = newID
	d.Path = newPath
	return d.save()
}

func (d *Document) save() error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(d.Path, data, 0644)
}
