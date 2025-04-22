[![CI](https://github.com/imhasandl/user-service/actions/workflows/ci.yml/badge.svg)](https://github.com/imhasandl/user-service/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/imhasandl/user-service)](https://goreportcard.com/report/github.com/imhasandl/user-service)
[![GoDoc](https://godoc.org/github.com/imhasandl/user-service?status.svg)](https://godoc.org/github.com/imhasandl/user-service)
[![Coverage](https://codecov.io/gh/imhasandl/user-service/branch/main/graph/badge.svg)](https://codecov.io/gh/imhasandl/user-service)
[![Go Version](https://img.shields.io/github/go-mod/go-version/imhasandl/user-service)](https://golang.org/doc/devel/release.html)

# User Service

A microservice for user management in a social media application, built with Go and gRPC.

## Overview

The User Service is responsible for managing user accounts, authentication, and profile information for the social media platform. It provides core functionality such as user registration, login, profile updates, and user data retrieval. The service uses gRPC for communication with other services in the microservices architecture.

## Prerequisites

- Go 1.23 or later
- PostgreSQL database
- RabbitMQ (for event-driven communication with other services)

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
PORT=":50053"
DB_URL="postgres://username:password@host:port/database?sslmode=disable"
# DB_URL="postgres://username:password@db:port/database?sslmode=disable" // FOR DOCKER COMPOSE
EMAIL="EMAIL_FOR_SENDING_EMAILS"
EMAIL_SECRET="SECRET_FOR_CONFIRMING_SENDING_EMAILS"
TOKEN_SECRET="YOUR_JWT_SECRET_KEY"
RABBITMQ_URL="amqp://username:password@host:port"
```

## Database Migrations

This service uses Goose for database migrations:

```bash
# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations postgres "YOUR_DB_CONNECTION_STRING" up
```

## gRPC Methods

The service implements the following gRPC methods:

### RegisterUser

Registers a new user in the system.

#### Request Format

```json
{
   "username": "johndoe",
   "email": "john.doe@example.com",
   "password": "securepassword",
   "full_name": "John Doe"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the created user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
   },
   "token": "JWT authentication token"
}
```

### LoginUser

Authenticates a user and returns a JWT token.

#### Request Format

```json
{
   "username_or_email": "johndoe OR john.doe@example.com",
   "password": "securepassword"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
   },
   "token": "JWT authentication token"
}
```

### GetUserProfile

Retrieves a user's profile information.

#### Request Format

```json
{
   "user_id": "UUID of the user"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "bio": "User biography text",
      "profile_picture_url": "URL to profile picture",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
   }
}
```

### UpdateUserProfile

Updates a user's profile information.

#### Request Format

```json
{
   "user_id": "UUID of the user",
   "full_name": "Updated Name",
   "bio": "Updated biography",
   "profile_picture_url": "URL to new profile picture"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "full_name": "Updated Name",
      "bio": "Updated biography",
      "profile_picture_url": "URL to new profile picture",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
   }
}
```

### ChangePassword

Changes a user's password.

#### Request Format

```json
{
   "user_id": "UUID of the user",
   "current_password": "currentSecurePassword",
   "new_password": "newSecurePassword"
}
```

#### Response Format

```json
{
   "success": true,
   "message": "Password successfully changed"
}
```

### ValidateToken

Validates a JWT token and returns the associated user information.

#### Request Format

```json
{
   "token": "JWT token string"
}
```

#### Response Format

```json
{
   "valid": true,
   "user_id": "UUID of the user if token is valid",
   "username": "Username if token is valid",
   "claims": {
      "additional": "claims",
      "from": "token"
   }
}
```

### GetUserByID

Retrieves a user's information by their ID.

#### Request Format

```json
{
   "id": "UUID of the user"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "subscribers": "Comma-separated list of subscriber IDs",
      "subscribed_to": "Comma-separated list of user IDs the user is subscribed to",
      "is_premium": false,
      "verification_code": 123456,
      "is_verified": true
   }
}
```

### GetUserByEmailOrUsername

Retrieves a user's information by their email or username.

#### Request Format

```json
{
   "identifier": "johndoe OR john.doe@example.com"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "subscribers": "Comma-separated list of subscriber IDs",
      "subscribed_to": "Comma-separated list of user IDs the user is subscribed to",
      "is_premium": false,
      "verification_code": 123456,
      "is_verified": true
   }
}
```

### GetUserByToken

Validates a JWT token and returns the associated user information.

#### Request Format

```json
{
   "token": "JWT token string"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "subscribers": "Comma-separated list of subscriber IDs",
      "subscribed_to": "Comma-separated list of user IDs the user is subscribed to",
      "is_premium": false,
      "verification_code": 123456,
      "is_verified": true
   }
}
```

### GetAllUsers

Retrieves information for all users in the system.

#### Request Format

```json
{}
```

#### Response Format

```json
{
   "users": [
      {
         "id": "UUID of the user 1",
         "username": "johndoe",
         "email": "john.doe@example.com",
         "created_at": "2023-01-01T12:00:00Z",
         "updated_at": "2023-01-01T12:00:00Z",
         "subscribers": "Comma-separated list of subscriber IDs",
         "subscribed_to": "Comma-separated list of user IDs the user is subscribed to",
         "is_premium": false,
         "verification_code": 123456,
         "is_verified": true
      },
      {
         "id": "UUID of the user 2",
         // Other user properties
      }
   ]
}
```

### ChangeUsername

Changes a user's username.

#### Request Format

```json
{
   "username": "newUsername"
}
```

#### Response Format

```json
{
   "user": {
      "id": "UUID of the user",
      "username": "newUsername",
      "email": "john.doe@example.com",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "subscribers": "Comma-separated list of subscriber IDs",
      "subscribed_to": "Comma-separated list of user IDs the user is subscribed to",
      "is_premium": false,
      "verification_code": 123456,
      "is_verified": true
   }
}
```

### ChangePassword

Changes a user's password.

#### Request Format

```json
{
   "password": "newSecurePassword"
}
```

#### Response Format

```json
{
   "status": "Password successfully changed"
}
```

### SubscribeUser

Subscribes a user to another user.

#### Request Format

```json
{
   "user_id": "UUID of the user to subscribe to"
}
```

#### Response Format

```json
{
   "status": true
}
```

### UnsubscribeUser

Unsubscribes a user from another user.

#### Request Format

```json
{
   "user_id": "UUID of the user to unsubscribe from"
}
```

#### Response Format

```json
{
   "status": true
}
```

### SendVerificationCode

Sends a verification code to a user's email.

#### Request Format

```json
{
   // Request details would depend on implementation
}
```

#### Response Format

```json
{
   "status": "Verification code sent successfully"
}
```

### ResetPassword

Resets a user's password using a verification code.

#### Request Format

```json
{
   "new_password": "newSecurePassword",
   "verification_code": 123456
}
```

#### Response Format

```json
{
   "status": "Password reset successfully"
}
```

### DeleteUser

Deletes a user account.

#### Request Format

```json
{
   "password": "currentPassword",
   "verify_message": "I understand this action cannot be undone"
}
```

#### Response Format

```json
{
   "status": "User account deleted successfully"
}
```

### DeleteAllUsers

Deletes all user accounts (admin function).

#### Request Format

```json
{}
```

#### Response Format

```json
{
   "status": "All user accounts deleted successfully"
}
```

## RabbitMQ Integration

The User Service publishes events to RabbitMQ when significant user actions occur, enabling other services to react accordingly.

### Event Publication

The service publishes events to:
- **Exchange**: `users.topic` (topic exchange)
- **Routing Keys**:
  - `user.registered` - When a new user registers
  - `user.updated` - When a user updates their profile
  - `user.deleted` - When a user account is deleted

### Message Format Example

```json
{
   "event_type": "user.registered",
   "user_id": "UUID of the user",
   "username": "johndoe",
   "timestamp": "2023-01-01T12:00:00Z",
   "data": {
      "email": "john.doe@example.com",
      "full_name": "John Doe"
   }
}
```

## Running the Service

```bash
go run cmd/main.go
```

## Docker Support

The service can be run using Docker:

```bash
# Build the Docker image
docker build -t user-service .

# Run the container
docker run -p 50053:50053 user-service
```

When deploying to different CPU architectures:

```bash
docker build --platform=linux/amd64 -t user-service .
```

For more details on Docker deployment, see [README.Docker.md](./README.Docker.md).


