// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: scraper.sql

package database

import (
	"context"
)

const getNextSongsToFetch = `-- name: GetNextSongsToFetch :many
SELECT id, name, artist_name, album_name, album_art, duration, video_id, is_available, release_date, created_at, updated_at FROM songs
WHERE is_available = FALSE
LIMIT $1
`

func (q *Queries) GetNextSongsToFetch(ctx context.Context, limit int32) ([]Song, error) {
	rows, err := q.db.QueryContext(ctx, getNextSongsToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Song
	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ArtistName,
			&i.AlbumName,
			&i.AlbumArt,
			&i.Duration,
			&i.VideoID,
			&i.IsAvailable,
			&i.ReleaseDate,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markSongAsAvailable = `-- name: MarkSongAsAvailable :exec
UPDATE songs
SET is_available = TRUE, updated_at = now()
WHERE id = $1
`

func (q *Queries) MarkSongAsAvailable(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, markSongAsAvailable, id)
	return err
}
