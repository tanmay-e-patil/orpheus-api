-- name: GetNextSongsToFetch :many
SELECT * FROM songs
WHERE is_available = FALSE
LIMIT $1;

-- name: MarkSongAsAvailable :exec
UPDATE songs
SET is_available = TRUE, updated_at = now()
WHERE id = $1;