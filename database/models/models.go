package models

import "sync"

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
	Data map[string]interface{} `json:"data"` // Key-value data
	Path string                 `json:"path"` // Path to the file on disk (optional)
}

// KeyValue represents a single key-value pair
type KeyValue struct {
	Key   string      `json:"key"`   // Key name
	Value interface{} `json:"value"` // Value associated with the key
}

// Response represents a generic response for operations
type Response struct {
	Success bool        `json:"success"` // Indicates success or failure
	Message string      `json:"message"` // Descriptive message
	Data    interface{} `json:"data"`    // Any additional data (optional)
}
