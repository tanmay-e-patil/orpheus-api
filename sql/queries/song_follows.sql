-- name: CreateSongFollow :one
INSERT INTO song_follows (id, song_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;