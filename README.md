
# Orpheus API

Orpheus API is a RESTful API for managing users and songs. It provides endpoints for user authentication, song creation, and retrieval.

## Prerequisites

- Go 1.18+
- PostgreSQL
- [godotenv](https://github.com/joho/godotenv)
- [pq](https://github.com/lib/pq)

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/tanmay-e-patil/orpheus-api.git
    cd orpheus-api
    ```

2. Create a `.env` file in the root directory with the following environment variables:
    ```env
    PORT=your_port
    DATABASE_URL=your_database_url
    JWT_SECRET=your_jwt_secret
    ```

3. Install dependencies:
    ```sh
    go mod tidy
    ```

4. Run the server:
    ```sh
    go run main.go
    ```

## Endpoints

### User Endpoints

- `POST /v1/users`: Sign up a new user.
- `POST /v1/users/login`: Sign in an existing user.
- `POST /v1/refresh`: Refresh the access token.
- `GET /v1/users/me`: Get the authenticated user's details.

### Song Endpoints

- `POST /v1/songs`: Create a new song (requires authentication).
- `GET /v1/songs`: Get all songs for the authenticated user.
- `GET /v1/songs/{songID}`: Get a song by its ID.

## Middleware

- `middlewareAuth`: Middleware to protect routes that require authentication.

## License

This project is licensed under the MIT License.