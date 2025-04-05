package documents

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"Build-your-own-database/database/models"
)

type DocumentManager struct {
	collection *models.Collection
	docMux     sync.RWMutex
}

// Constructor
func NewDocumentManager(collection *models.Collection) *DocumentManager {
	return &DocumentManager{
		collection: collection,
	}
}

// Generate a random ID (internal use only)
func generateRandomID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 1. CreateDocument (by name)
func (dm *DocumentManager) CreateDocument(name string, data map[string]interface{}) (*models.Document, error) {
	dm.docMux.Lock()
	defer dm.docMux.Unlock()

	// Check if name already exists
	for _, doc := range dm.collection.Documents {
		if doc.Name == name {
			return nil, fmt.Errorf("document with name '%s' already exists", name)
		}
	}

	id := generateRandomID()
	docPath := filepath.Join(dm.collection.Path, id+".json")
	doc := &models.Document{
		ID:   id,
		Name: name,
		Data: data,
		Path: docPath,
	}

	dm.collection.Documents[id] = doc

	// Save to disk
	file, err := os.Create(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create document file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(doc); err != nil {
		return nil, fmt.Errorf("failed to encode document: %v", err)
	}

	fmt.Println("Created document:", name)
	return doc, nil
}

// 2. UseDocument (by name)
func (dm *DocumentManager) UseDocument(name string) (*models.Document, error) {
	dm.docMux.RLock()
	for _, doc := range dm.collection.Documents {
		if doc.Name == name {
			dm.docMux.RUnlock()
			fmt.Println("Using document from memory:", name)
			return doc, nil
		}
	}
	dm.docMux.RUnlock()

	// Not in memory? Load from disk
	files, err := os.ReadDir(dm.collection.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(dm.collection.Path, file.Name())
			f, err := os.Open(path)
			if err != nil {
				continue
			}

			var doc models.Document
			if err := json.NewDecoder(f).Decode(&doc); err == nil {
				f.Close()
				if doc.Name == name {
					doc.Path = path
					dm.docMux.Lock()
					dm.collection.Documents[doc.ID] = &doc
					dm.docMux.Unlock()
					fmt.Println("Loaded document from disk:", name)
					return &doc, nil
				}
			}
			f.Close()
		}
	}

	return nil, fmt.Errorf("document '%s' does not exist", name)
}

// 3. DeleteDocument (by name)
func (dm *DocumentManager) DeleteDocument(name string) error {
	dm.docMux.Lock()
	defer dm.docMux.Unlock()

	for id, doc := range dm.collection.Documents {
		if doc.Name == name {
			if err := os.Remove(doc.Path); err != nil {
				return fmt.Errorf("failed to delete document file: %v", err)
			}
			delete(dm.collection.Documents, id)
			fmt.Println("Deleted document:", name)
			return nil
		}
	}

	return fmt.Errorf("document '%s' does not exist", name)
}

// 4. RenameDocument (by name)
func (dm *DocumentManager) RenameDocument(oldName, newName string) error {
	dm.docMux.Lock()
	defer dm.docMux.Unlock()

	// Check if newName already exists
	for _, d := range dm.collection.Documents {
		if d.Name == newName {
			return fmt.Errorf("document '%s' already exists", newName)
		}
	}

	for _, doc := range dm.collection.Documents {
		if doc.Name == oldName {
			doc.Name = newName

			// Save with updated name
			file, err := os.Create(doc.Path)
			if err != nil {
				return fmt.Errorf("failed to update renamed doc: %v", err)
			}
			defer file.Close()
			if err := json.NewEncoder(file).Encode(doc); err != nil {
				return fmt.Errorf("failed to encode renamed doc: %v", err)
			}

			fmt.Printf("Renamed document '%s' to '%s'\n", oldName, newName)
			return nil
		}
	}

	return fmt.Errorf("document '%s' not found", oldName)
}

// 5. FindDocument (by key-value inside data)
func (dm *DocumentManager) FindDocument(key string, val interface{}) []*models.Document {
	dm.docMux.RLock()
	defer dm.docMux.RUnlock()

	var results []*models.Document
	for _, doc := range dm.collection.Documents {
		if v, ok := doc.Data[key]; ok && v == val {
			results = append(results, doc)
		}
	}

	fmt.Printf("Found %d document(s) matching %s = %v\n", len(results), key, val)
	return results
}
