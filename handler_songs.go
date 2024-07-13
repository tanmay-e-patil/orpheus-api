package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
)

const assetsDir = "./public/library"

func (cfg *apiConfig) handlerSongsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type Song struct {
		SpotifySongID          string `json:"spotify_song_id"`
		SpotifySongName        string `json:"spotify_song_name"`
		SpotifySongArtist      string `json:"spotify_song_artist"`
		SpotifySongAlbum       string `json:"spotify_song_album"`
		YtSongDuration         string `json:"yt_song_duration"`
		SpotifySongAlbumArtURL string `json:"spotify_song_album_art_url"`
		SpotifyReleaseDate     string `json:"spotify_release_date"`
		YtSongVideoID          string `json:"yt_song_video_id"`
	}
	type parameters struct {
		Song Song `json:"song"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	song, err := cfg.DB.GetSongById(r.Context(), params.Song.SpotifySongID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	} else {
		log.Printf("Song already exists. No need to create")
		songFollow, err := cfg.DB.CreateSongFollow(r.Context(), database.CreateSongFollowParams{
			ID:     id,
			SongID: params.Song.SpotifySongID,
			UserID: user.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		log.Printf("Created song follow: %+v", songFollow)
		respondWithJSON(w, http.StatusCreated, databaseSongToSong(song))
	}
	date, err := time.Parse("2006-01-02", params.Song.SpotifyReleaseDate)
	if err != nil {
		log.Printf("Couldn't parse spotify release date: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Couldn't parse spotify release date")
		return
	}

	song, err = cfg.DB.CreateSong(r.Context(), database.CreateSongParams{
		ID:          params.Song.SpotifySongID,
		Name:        params.Song.SpotifySongName,
		ArtistName:  params.Song.SpotifySongArtist,
		AlbumName:   params.Song.SpotifySongAlbum,
		AlbumArt:    params.Song.SpotifySongAlbumArtURL,
		Duration:    params.Song.YtSongDuration,
		VideoID:     params.Song.YtSongVideoID,
		ReleaseDate: date,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create song")
	}

	songFollow, err := cfg.DB.CreateSongFollow(r.Context(), database.CreateSongFollowParams{
		ID:     id,
		SongID: params.Song.SpotifySongID,
		UserID: user.ID,
	})

	log.Printf("Created song follow: %+v", songFollow)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, databaseSongToSong(song))
}

func (cfg *apiConfig) handlerSongsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	songs, err := cfg.DB.GetSongsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, databaseSongsToSongs(songs))
}

func (cfg *apiConfig) handlerSongsGetByID(w http.ResponseWriter, r *http.Request, user database.User) {
	songID := r.PathValue("songID")
	song, err := cfg.DB.GetSongById(r.Context(), songID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	songFilePath := fmt.Sprintf("%s/%s/%s/%s.m4a", assetsDir, song.ArtistName, song.AlbumName, song.Name)

	// Set the content-type header for M4A files
	w.Header().Set("Content-Type", "audio/mp4")

	// Serve the M4A file
	http.ServeFile(w, r, songFilePath)
}
