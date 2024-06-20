-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, password)
VALUES ($1, $2, $3, $4,  $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;