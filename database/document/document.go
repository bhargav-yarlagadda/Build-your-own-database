package documents

import (
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

// 1. CreateDocument
func (dm *DocumentManager) CreateDocument(id string, data map[string]interface{}) (*models.Document, error) {
	dm.docMux.Lock()
	defer dm.docMux.Unlock()

	if _, exists := dm.collection.Documents[id]; exists {
		return nil, fmt.Errorf("document '%s' already exists", id)
	}

	docPath := filepath.Join(dm.collection.Path, id+".json")
	doc := &models.Document{
		ID:   id,
		Data: data,
		Path: docPath,
	}

	// Save to memory
	dm.collection.Documents[id] = doc

	// Save to disk
	file, err := os.Create(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create document file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(doc); err != nil {
		return nil, fmt.Errorf("failed to encode document: %v", err)
	}

	fmt.Println("Created document:", id)
	return doc, nil
}

// 2. UseDocument (load from disk if not in memory)
func (dm *DocumentManager) UseDocument(id string) (*models.Document, error) {
	dm.docMux.RLock()
	doc, exists := dm.collection.Documents[id]
	dm.docMux.RUnlock()

	if exists {
		fmt.Println("Using document from memory:", id)
		return doc, nil
	}

	// Load from disk
	docPath := filepath.Join(dm.collection.Path, id+".json")
	file, err := os.Open(docPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("document '%s' does not exist", id)
		}
		return nil, err
	}
	defer file.Close()

	var loadedDoc models.Document
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedDoc); err != nil {
		return nil, err
	}
	loadedDoc.Path = docPath

	// Save to memory
	dm.docMux.Lock()
	dm.collection.Documents[id] = &loadedDoc
	dm.docMux.Unlock()

	fmt.Println("Loaded document from disk:", id)
	return &loadedDoc, nil
}

// 3. DeleteDocument
func (dm *DocumentManager) DeleteDocument(id string) error {
	dm.docMux.Lock()
	defer dm.docMux.Unlock()

	doc, exists := dm.collection.Documents[id]
	if !exists {
		return fmt.Errorf("document '%s' does not exist", id)
	}

	// Delete from disk
	if err := os.Remove(doc.Path); err != nil {
		return fmt.Errorf("failed to delete document file: %v", err)
	}

	// Delete from memory
	delete(dm.collection.Documents, id)

	fmt.Println("Deleted document:", id)
	return nil
}

// 4. FindDocument
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
