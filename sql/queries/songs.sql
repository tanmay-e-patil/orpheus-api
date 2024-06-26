-- name: CreateSong :one
INSERT INTO songs (id, name, artist_name, album_name, album_art, duration, video_id, created_at, updated_at)
VALUES ($1, $2, $3, $4,  $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;