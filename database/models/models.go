package models

// Database represents a database in your system
type Database struct {
	Name      string     `json:"name"`      // Name of the database
	Documents []Document `json:"documents"` // List of documents in the database
}

// Document represents a document stored in a database
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
