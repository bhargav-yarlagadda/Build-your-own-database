


# Build Your Own Database 
🚀 I'm currently working on the `db-revamp` branch.  
If you'd like to contribute, pull the latest changes from `db-revamp` for the latest updates or `dev` for a more stable experience.  
Submit a PR, and let's make this repo even better! 😃  
Happy coding! 🚀  

This is a simple database implementation using Go and Fiber. It provides an API for managing databases and documents, allowing users to create, delete, and manage documents within databases. This project is built with Go's powerful standard library and the Fiber web framework for creating RESTful APIs.

## Features

- Create a database
- Delete a database
- Create a document within a database
- Delete a document from a database
- access attributes in the documents and modify  them.
- Implemented concurrency control mechanisms to manage shared resource locking during simultaneous database edits, ensuring data consistency and preventing race conditions
## Technologies Used

- **Go**: The programming language for building the database and API.
- **Fiber**: A fast and lightweight web framework for Go.
- **UUID**: Used for generating unique identifiers for documents.
- **JSON**: Used for document storage.



## Project Structure

```
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
```


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

## API Endpoints

# API Endpoints

- **POST** `/create-database`
- **DELETE** `/delete-database/{databaseName}`
- **POST** `/{databaseName}/create-document`
- **DELETE** `/{databaseName}/delete-document/{documentName}`
- **GET** `/{databaseName}/documents`
- **GET** `/{databaseName}/document/{documentName}`
- **PATCH** `/{databaseName}/update-document/{documentName}`
- **PATCH** `/{databaseName}/{documentName}/delete-pair/{key}`
- **PATCH** `/{databaseName}/{documentName}/update-pair`
## ToDo's
- Distributed File Storage
