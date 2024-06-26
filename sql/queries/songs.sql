-- name: CreateSong :one
INSERT INTO songs (id, name, artist_name, album_name, album_art, duration, video_id)
VALUES ($1, $2, $3, $4,  $5, $6, $7)
RETURNING *;