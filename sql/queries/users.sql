-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, password, refresh_token)
VALUES ($1, $2, $3, $4,  $5, $6, $7)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * from users
WHERE id = $1;

-- name: UpdateUserWithRefreshToken :one
UPDATE users SET refresh_token = $1
WHERE id = $2
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT * from users
WHERE refresh_token = $1;