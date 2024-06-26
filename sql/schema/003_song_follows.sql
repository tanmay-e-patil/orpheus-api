-- +goose Up
CREATE TABLE song_follows (
    id UUID PRIMARY KEY NOT NULL,
    song_id TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_songs FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE song_follows;