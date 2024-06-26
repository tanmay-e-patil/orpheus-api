package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerSongsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type Song struct {
		SpotifySongID          string `json:"spotify_song_id"`
		SpotifySongName        string `json:"spotify_song_name"`
		SpotifySongArtist      string `json:"spotify_song_artist"`
		SpotifySongAlbum       string `json:"spotify_song_album"`
		YtSongDuration         string `json:"yt_song_duration"`
		SpotifySongAlbumArtURL string `json:"spotify_song_album_art_url"`
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

	song, err := cfg.DB.GetSongById(r.Context(), params.Song.SpotifySongID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	} else {
		log.Printf("Song already exists. No need to create")
		respondWithJSON(w, http.StatusCreated, databaseSongToSong(song))
	}

	song, err = cfg.DB.CreateSong(r.Context(), database.CreateSongParams{
		ID:         params.Song.SpotifySongID,
		Name:       params.Song.SpotifySongName,
		ArtistName: params.Song.SpotifySongArtist,
		AlbumName:  params.Song.SpotifySongAlbum,
		AlbumArt:   params.Song.SpotifySongAlbumArtURL,
		Duration:   params.Song.YtSongDuration,
		VideoID:    params.Song.YtSongVideoID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create song")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
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
