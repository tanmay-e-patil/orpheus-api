-- name: CreateSong :one
INSERT INTO songs (id, name, artist_name, album_name, album_art, duration, video_id, release_date, created_at, updated_at)
VALUES ($1, $2, $3, $4,  $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;

-- name: GetSongsForUser :many
SELECT * from songs
WHERE id IN (
    SELECT song_id FROM song_follows
                   WHERE user_id = $1
    );