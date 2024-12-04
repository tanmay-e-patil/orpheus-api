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

#### Sign Up

- **URL:** `POST /v1/users`
- **Description:** Sign up a new user.
- **Request Body:**
    ```json
    {
        "username": "string",
        "email": "string",
        "password": "string"
    }
    ```
- **Response:**
    ```json
    {
        "id": "string",
        "username": "string",
        "email": "string",
        "access_token": "string"
    }
    ```

#### Sign In

- **URL:** `POST /v1/users/login`
- **Description:** Sign in an existing user.
- **Request Body:**
    ```json
    {
        "email": "string",
        "password": "string"
    }
    ```
- **Response:**
    ```json
    {
        "id": "string",
        "username": "string",
        "email": "string",
        "access_token": "string"
    }
    ```

#### Refresh Token

- **URL:** `POST /v1/refresh`
- **Description:** Refresh the access token.
- **Response:**
    ```json
    {
        "access_token": "string"
    }
    ```

#### Get User Details

- **URL:** `GET /v1/users/me`
- **Description:** Get the authenticated user's details.
- **Response:**
    ```json
    {
        "id": "string",
        "username": "string",
        "email": "string"
    }
    ```

### Song Endpoints

#### Create Song

- **URL:** `POST /v1/songs`
- **Description:** Create a new song (requires authentication).
- **Request Body:**
    ```json
    {
        "song": {
            "spotify_song_id": "string",
            "spotify_song_name": "string",
            "spotify_song_artist": "string",
            "spotify_song_album": "string",
            "yt_song_duration": "string",
            "spotify_song_album_art_url": "string",
            "spotify_release_date": "string",
            "yt_song_video_id": "string"
        }
    }
    ```
- **Response:**
    ```json
    {
        "id": "string",
        "name": "string",
        "artist_name": "string",
        "album_name": "string",
        "album_art": "string",
        "duration": "string",
        "video_id": "string",
        "release_date": "string",
        "created_at": "string",
        "updated_at": "string"
    }
    ```

#### Get All Songs

- **URL:** `GET /v1/songs`
- **Description:** Get all songs for the authenticated user.
- **Response:**
    ```json
    [
        {
            "id": "string",
            "name": "string",
            "artist_name": "string",
            "album_name": "string",
            "album_art": "string",
            "duration": "string",
            "video_id": "string",
            "release_date": "string",
            "created_at": "string",
            "updated_at": "string"
        }
    ]
    ```

#### Get Song by ID

- **URL:** `GET /v1/songs/{songID}`
- **Description:** Get a song by its ID.
- **Response:**
    ```json
    {
        "id": "string",
        "name": "string",
        "artist_name": "string",
        "album_name": "string",
        "album_art": "string",
        "duration": "string",
        "video_id": "string",
        "release_date": "string",
        "created_at": "string",
        "updated_at": "string"
    }
    ```

## Middleware

- `middlewareAuth`: Middleware to protect routes that require authentication.

## License

This project is licensed under the MIT License.