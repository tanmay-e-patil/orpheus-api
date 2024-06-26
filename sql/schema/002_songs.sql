CREATE TABLE songs(
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    artist_name TEXT NOT NULL,
    album_name TEXT NOT NULL,
    album_art TEXT NOT NULL,
    duration TEXT NOT NULL,
    video_id TEXT NOT NULL
);