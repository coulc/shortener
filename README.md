## Shortener
A simple URL shortener API service built with Go + SQLite.

## Features
* Create short links   
* Redirect to original links   
* Delete short links   
* Visit count tracking  
* slog-based logging system (supports both file and console output)  
* SQLite storage with automatic table creation  
* Unit tests included  

## Tech Stack
Go v1.25.5  
SQLite3 v1.14.42A  
net/http - HTTP server and client  
log/slog - Structured logging  
mattn/go-sqlite3 - SQLite driver  

## Project Structure
```text
.
├── cmd/                       # Application entry points
│   ├── main.go                # Main entry point
│   └── setup.go               # Initialization helpers
├── db/                        # Database directory
│   └── urls.db                # SQLite database file
├── handler/                   # HTTP handlers
│   └── shortener.go           # Create, Redirect, Delete handlers
├── storage/                   # Data storage layer
│   ├── storage.go             # Storage interface
│   └── sqlite.go              # SQLite implementation
├── test/                      # Unit tests
│   └── shortener_test.go      # Handler tests
├── utils/                     # Utilities
│   ├── hash.go                # Short code generator
│   └── logger.go              # Logging initialization
├── log/                       # Log directory
│   └── logs.log               # Application logs
├── go.mod
├── go.sum
└── README.md

```
## Getting Started
### Prerequisites
* Go 1.21 or higher
* SQLite3 (no separate installation required, pure Go driver)

### Installation
```bash
git clone https://github.com/coulc/shortener
cd shortener
go mod tidy
```
### Run the Service
```bash
go run cmd/*.go
The server will start on http://localhost:8080
```

### Run Tests
```bash
go test ./test -v
```

## API Documentation
### Base URL
```text
http://localhost:8080
```

### Endpoints
```bash
Method	Endpoint	Description	Status Codes
POST	/	Create a new short link	201, 409, 400
GET	/{shortCode}	Redirect to original URL	302, 404, 400
DELETE	/{shortCode}	Delete a short link	204, 404, 400
```

### 1. Create Short Link
Create a shortened URL from a long URL.
Request:
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"long_url":"https://github.com/coulc"}'

```
Response(201 Created):
```json
{
  "short_code": "2ShiN0",
  "short_url": "http://localhost:8080/2ShiN0",
  "long_url": "https://github.com/coulc"
}
```

### 2. Redirect to Original URL
Access the short link to be redirected to the original URL.
Request:
```bash
curl -v http://localhost:8080/2ShiN0
```

### 3. Delete Short Link
Delete an existing short link.
Request:

```bash
curl -X DELETE http://localhost:8080/2ShiN0
```
Response: 
```text
204 No Content
```

## Database Schema
```sql
CREATE TABLE IF NOT EXISTS urls (
    short_code TEXT PRIMARY KEY,
    long_url TEXT NOT NULL UNIQUE,
    created_at INTEGER,
    visit_count INTEGER
);
```

## Short Code Generation  
Short codes are generated using a hash function based on the long URL. The same URL will always generate the same short code, ensuring consistency.

## Error Handling
```text
HTTP Status	Description	When it occurs
200	OK	Successful operations
201	Created	Short link successfully created
204	No Content	Successful deletion
302	Found	Redirect to original URL
400	Bad Request	Invalid URL format or missing parameters
404	Not Found	Short code doesn't exist
409	Conflict	URL already exists in system
500	Internal Server Error	Database or server errors
```
