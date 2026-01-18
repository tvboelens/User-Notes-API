# User-Notes-API
## Description and project motivation
A production-like Go backend service for managing user notes, with JWT authentication, PostgreSQL/MySQL support, and RESTful endpoints.

This project was built as a hands-on exercise to learn Go and its ecosystem,
including HTTP servers (Gin), authentication (JWT), cryptography (Argon2id),
and database integration (GORM).

## Features
- User registration and login with JWT authentication
- Password hashing using Argon2id
- CRUD operations for notes (create, read, update, delete)
- REST API implemented with Gin
- Database integration using GORM (Postgres or MySQL)
- Configuration via .env
- Unit tests for services, controllers, and middleware
## Project structure and design
The server is implemented using a controller-services-repositories pattern. Gin routes an HTTP request to the corresponding controller, which handles the HTTP request itself. The controller then calls the relevant service to obtain the data for the response. At this moment there are two services:
- Auth service, which handles business logic of registration and login.
- Note service, which handles business logic of CRUD operations on notes.

The services access the database via the repositories, which themselves use GORM.

Authentication using JWT tokens is implemented through Gin middleware.
## Getting started
### Prerequisites
- Go 1.25+
- Docker  & Docker Compose (for local DB)
- Postgres (if not using Docker)
### Available endpoints
| Method | Path | Auth | Description
|--------|------|------|------------|
|POST | `/register` | No | Register new user
|POST | `/login` | No | Login with username and password
| POST | `/notes` | Yes | Create new note
| GET | `/notes` | Yes | Get the ids and titles of all notes belonging to specific user
| GET | `/notes/:id` | Yes | Get note with a specific id
| UPDATE | `/notes/:id` | Yes | update a note
| DELETE | `/notes/:id` | Yes | delete a note

**Authorization:** Include header:
```
Authorization: Bearer <your_jwt_token>
```

### Running tests
**Unit tests:**
```
go test ./... --short
```
**Integration tests:**
```
go test ./testing/integration
```
**End-to-end-test:** Run the script in `scripts/run-e2e.sh`. This requires an installation of Docker Compose.
