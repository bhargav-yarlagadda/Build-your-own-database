package models

// GoDB is the central database manager
type GoDB struct {
	Databases map[string]*Database // Stores all databases
}

// Database represents a database in the system
type Database struct {
	Name        string                  `json:"name"`        // Name of the database
	Collections map[string]*Collection  `json:"collections"` // List of collections in the database
}

// Collection represents a collection inside a database
type Collection struct {
	Name      string               `json:"name"`      // Name of the collection
	Documents map[string]*Document `json:"documents"` // List of documents in the collection
}

// Document represents a document stored in a collection
type Document struct {
	UUID    string                 `json:"uuid"`    // Unique identifier for the document
	Name    string                 `json:"name"`    // Name of the document
	Content map[string]interface{} `json:"content"` // Key-value pairs in the document
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
