# Build Your Own Database

This is a simple database implementation using Go and Fiber. It provides an API for managing databases and documents, allowing users to create, delete, and manage documents within databases. This project is built with Go's powerful standard library and the Fiber web framework for creating RESTful APIs.

## Features

- Create a database
- Delete a database
- Create a document within a database
- Delete a document from a database
- access attributes in the documents and modify  them.
## Technologies Used

- **Go**: The programming language for building the database and API.
- **Fiber**: A fast and lightweight web framework for Go.
- **UUID**: Used for generating unique identifiers for documents.
- **JSON**: Used for document storage.

## Project Structure

├── api/
│   ├── handlers/
│   │   ├── create_document.go   # Handler for creating documents
│   │   └── delete_document.go   # Handler for deleting documents
│   ├── routes/
│   │   └── routes.go            # Defines all API routes
├── config/config.go
├── database/
│   ├── db/
│   │   └── database.go          # Contains database-related functions
│   ├── document/
│   │   └── document.go          # Functions for document creation and management
│   ├── key-valueskey-values/
│   │   └── document.go          # Functions for accessing attributes in doucments
│   └── models/
│       └── models.go          # Data models (e.g., document)
│   └── utils/
│       └── utils.go          # Utility functions 
├── main.go                      # Entry point for the Go server
├── go.mod                      
├── go.sum
└── README.md                    # MardDown


## installation
```bash
  git clone https://github.com/yourusername/build-your-own-database.git
  cd build-your-own-database

```
install the dependencies
```bash
go mod tidy 
```

start the api
```bash
go run main.go
```

### navigate to config.go and change the basePath as per your preference

## api end points
## API Endpoints

### 1. **Create a Database**
- **URL**: `/create-database`
- **Method**: `POST`

### 2. **Delete a Database**
- **URL**: `/delete-database/{databaseName}`
- **Method**: `DELETE`

### 3. **Create a Document**
- **URL**: `/{databaseName}/create-document`
- **Method**: `POST`

### 4. **Delete a Document**
- **URL**: `/{databaseName}/delete-document/{documentName}`
- **Method**: `DELETE`

### 5. **Get All Documents in a Database**
- **URL**: `/{databaseName}/documents`
- **Method**: `GET`

### 6. **Get a Document by Name**
- **URL**: `/{databaseName}/document/{documentName}`
- **Method**: `GET`

