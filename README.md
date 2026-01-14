# User API

A simple REST API for managing user data built with Go. The API stores user information in memory as a JSON map and provides endpoints for creating, reading, updating, and deleting users.


## Run the server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### GET /api/users

Retrieve all users sorted alphabetically by name.

Response:
```json
{
  "data": [
    "Alice Smith id: 15",
    "John Doe id: 1"
  ]
}
```

### GET /api/users/{id}

Retrieve a specific user by ID.

Response:
```json
{
  "data": {
    "first_name": "Alice",
    "last_name": "Smith",
    "biography": "An artist"
  }
}
```

### POST /api/users

Create a new user. The ID is incremental and automatically assigned.

Request body:
```json
{
  "first_name": "Alice",
  "last_name": "Smith",
  "biography": "An artist"
}
```

Response:
```json
{
  "data": "New id = 16"
}
```

### PUT /api/users/{id}

Update an existing user by ID.

Request body:
```json
{
  "first_name": "Alice",
  "last_name": "Smith",
  "biography": "Updated biography"
}
```

Response:
```json
{
  "data": "User updated = {Alice Smith Updated biography}"
}
```

### DELETE /api/users/{id}

Delete a user by ID.

Response:
```json
{
  "data": "Deleted user id: 15"
}
```

## Data Validation

All user fields are required:

- `first_name` - cannot be empty
- `last_name` - cannot be empty
- `biography` - cannot be empty

If any field is missing, the API returns a 400 Bad Request error.

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful operation
- `400 Bad Request` - Invalid input or missing required fields
- `404 Not Found` - User not found
- `413 Payload Too Large` - Request body exceeds 1024 bytes
- `500 Internal Server Error` - Server error

## Server Configuration

- Port: 8080
- Max request body size: 1024 bytes
- Read timeout: 10 seconds
- Write timeout: 10 seconds
- Idle timeout: 1 minute

## Project Structure

- `user` - Struct for user data
- `application` - Struct containing the user map database
- `Response` - Struct for API responses
- `findAll()` - Get all users
- `findByID()` - Get user by ID
- `insertNewUser()` - Create new user
- `updateUser()` - Update existing user
- `deleteUser()` - Delete user
- Handler functions - HTTP request handlers for each endpoint
- `sendJSON()` - Helper to send JSON responses
