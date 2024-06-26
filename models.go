package main

import (
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type Song struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ArtistName string    `json:"artist_name"`
	AlbumName  string    `json:"album_name"`
	AlbumArt   string    `json:"album_art"`
	Duration   string    `json:"duration"`
	VideoID    string    `json:"video_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func databaseSongToSong(song database.Song) Song {
	return Song{
		ID:         song.ID,
		Name:       song.Name,
		ArtistName: song.ArtistName,
		AlbumName:  song.AlbumName,
		AlbumArt:   song.AlbumArt,
		Duration:   song.Duration,
		VideoID:    song.VideoID,
		CreatedAt:  song.CreatedAt,
		UpdatedAt:  song.UpdatedAt,
	}
}

func databaseUserToUser(user database.User, accessToken string) User {
	return User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	}
}
