

# Build Your Own Database  
📦 A lightweight, file-based database built from scratch using Go – **no external DBs**, **no APIs**, just pure Go code!  
Everything is handled via the `main.go` file, which acts as the client interface for interacting with the database system.

---

## 🧠 What Is This?

This project is a simple yet powerful key-value document database system implemented entirely in Go. It allows users to:

- Create/delete **databases**
- Create/delete **documents** inside databases
- **Access**, **modify**, **rename**, and **delete** key-value pairs inside documents
- **Fetch documents** using their **name**
- Perform all operations from the **command line**, with a focus on modular, extensible architecture

All data is stored **locally in files** and persists between sessions. Internally, the system uses **UUIDs**, but users interact with it using **names** for ease of access.

---

## ✨ Features

- ✅ Create & delete databases  
- ✅ Create & delete documents by name  
- ✅ Add/update/delete key-value pairs in a document  
- ✅ Fetch documents by name  
- ✅ Rename document names  
- ✅ File-based storage (JSON)  
- ✅ Concurrency-safe using Go mutexes  
- ✅ Modular code structure  

---

## 🛠 Technologies Used

- **Go** – Core language  
- **JSON** – For persistent document storage  
- **UUID** – For unique internal document IDs  
- **Mutex** – For handling concurrent operations safely  

---

## 📁 Project Structure

```
├── config/
│   └── config.go                  # Configuration for base file path
├── database/
│   ├── db/
│   │   └── database.go            # Functions for DB creation/deletion
│   ├── document/
│   │   └── document.go           # Document creation/deletion, renaming
│   │   └── document.go           # Add/update/delete key-value pairs
│   ├── models/
│   │   └── models.go             # Data models for DB and documents
│   └── utils/
│       └── utils.go              # File and helper utilities
├── main.go                        # CLI entry point for all operations
├── go.mod                         # Go module definition
├── go.sum                         # Go dependency checksum
└── README.md                      # You're reading it!
```

---

## 🚀 Getting Started

```bash
git clone https://github.com/bhargav-yarlagadda/build-your-own-database.git
cd build-your-own-database
```

Install dependencies:

```bash
go mod tidy
```

Run the project:
**main.go acts as a client that interacts with the db. take a glimpse of main.go to understand the working of the db.**
```bash
go run main.go
```

> **Note**: You can modify the default file storage path by updating `config/config.go`.

---
# Refactoring `dbManager.go` into `document_manager.go` and `collection_manager.go`

## Motivation
The original `dbManager.go` file was handling logic for both collections and documents in a single place. This violated the **Single Responsibility Principle** and made the codebase harder to maintain and extend.

## Changes Made

### 1. Created New Files
- `document_manager.go`: Handles all document-specific operations.
- `collection_manager.go`: Handles all collection-specific operations.

---

### 2. Document Refactor Highlights

#### ✅ New Struct: `DocumentManager`
- Encapsulates all document-specific CRUD operations.

#### ✅ Thread-Safety
- Introduced `sync.RWMutex` (`docMux`) to ensure concurrent read/write safety while accessing or modifying documents.

#### ✅ Clear Method Separation
Each method is focused on a single task:
- `CreateDocument`
- `UseDocument`
- `UpdateDocument`
- `DeleteDocument`
- `FetchDocument`
- `RenameDocument`
- `DeleteKey`

#### ✅ File I/O Improvements
- Used `json.NewEncoder`/`Decoder` consistently.
- Ensured proper file closing using `defer`.

#### ✅ Error Handling
- Added meaningful error messages.
- Ensured consistency in error format and logging.

#### ✅ In-Memory Caching
- When a document is used, it is loaded into memory if not already present.

---

### 3. Collection Refactor Highlights

#### ✅ New Struct: `CollectionManager`
- Handles collection-level logic such as:
  - `CreateCollection`
  - `DeleteCollection`
  - `ListCollections`
  - `RenameCollection`
  - `LoadCollection` from disk

#### ✅ Directory Structure
- Each collection has its own subdirectory inside `./data`.

#### ✅ Improved Initialization
- Clean separation between initializing a collection and working with documents inside it.

#### ✅ Mutex for Collection Safety
- Collection-level operations are also thread-safe using `sync.Mutex`.

---
## Benefits
- **Cleaner Code Structure**: Collections and documents now handled separately.
- **Easier to Maintain**: Logical grouping of responsibilities.
- **Better Concurrency**: Thread-safe reads and writes.
- **Improved Readability**: Self-explanatory function names and simplified logic.

## To Do
- Add unit tests for `DocumentManager` and `CollectionManager`
- Extend functionality with indexing or search support
- Integrate logging system instead of `fmt.Println`

## Example Usage
```go
collection := collectionManager.CreateCollection("users")

docManager := NewDocumentManager(collection)
docManager.CreateDocument("user1", map[string]interface{}{"name": "John"})
```

---
✅ Refactored, modular, and scalable!





---

## 📌 To-Do's

- Implement Distributed File Storage  
- Add CLI interface with flags (optional)  
- Enable nested key support  
- Add document versioning (optional history)  
- Unit tests for each module  

---

## 🙌 Contribute

Fork the repo, make your changes, and raise a PR!  
Suggestions and improvements are always welcome. 😊

---

Let me know if you'd like me to generate badges, screenshots, or usage gifs for a finishing touch!
