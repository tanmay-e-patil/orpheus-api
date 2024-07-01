package main

import (
	"context"
	"fmt"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

const BASE_DIR = "./public/library"

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		songs, err := db.GetNextSongsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next songs to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(songs))

		wg := &sync.WaitGroup{}
		for _, song := range songs {
			wg.Add(1)
			go scrapeSongs(db, wg, song)
		}
		wg.Wait()
	}
}

func createSongDirectories(dirPath string) error {
	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory along with any necessary parents
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
		fmt.Println("Directory created:", dirPath)
	} else {
		fmt.Println("Directory already exists:", dirPath)
	}
	return nil

}

func downloadCoverArt(song database.Song) (string, error) {
	// Create the file

	outputFilePath := fmt.Sprintf("%s/%s/%s/cover_art.jpeg", BASE_DIR, song.ArtistName, song.AlbumName)
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return "", err
	}
	defer file.Close()

	// Get the data
	resp, err := http.Get(song.AlbumArt)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status: %s\n", resp.Status)
		return "", err
	}

	// Writer the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return "", err
	}
	return outputFilePath, nil
}

func downloadSong(song database.Song) error {
	// Create appropriate directories for song using BASE_DIR/ArtistName/AlbumName/
	songDirPath := fmt.Sprintf("%s/%s/%s", BASE_DIR, song.ArtistName, song.AlbumName)
	err := createSongDirectories(songDirPath)
	if err != nil {
		log.Printf("Error creating song directories: %v", err)
		return err
	}
	// Download the song in that directory
	// example cli command: yt-dlp -f bestaudio --extract-audio --audio-format flac --audio-quality 0
	// xfqBQ2XhBCg -o SPOT! --postprocessor-args
	//"ffmpeg:-metadata title='SPOT!' -metadata artist='ZICO, JENNIE' -metadata album='SPOT!' -metadata date='2024' -metadata year='2024'"

	songFilePath := songDirPath + fmt.Sprintf("/tmpSong")

	songMetadata := fmt.Sprintf("ffmpeg:-metadata title='%s' -metadata artist='%s' -metadata album='%s' -metadata date='%v' -metadata year='%v'",
		song.Name, song.ArtistName, song.AlbumName, song.ReleaseDate.Year(), song.ReleaseDate.Year())

	ytDlpCmd := exec.Command("yt-dlp",
		"-f", "bestaudio",
		"--extract-audio",
		"--audio-format", "m4a",
		"--audio-quality", "0",
		fmt.Sprintf(song.VideoID),
		"-o", songFilePath,
		"--postprocessor-args",
		songMetadata,
	)

	ytDlpCmd.Stdout = os.Stdout
	ytDlpCmd.Stderr = os.Stderr

	log.Printf("Downloading and converting audio...")
	log.Printf("%v", ytDlpCmd.Args)
	err = ytDlpCmd.Run()
	if err != nil {
		fmt.Printf("Error running yt-dlp command: %v\n", err)
		return err
	}

	// Download the cover art
	coverArtPath, err := downloadCoverArt(song)
	if err != nil {
		log.Printf("Error downloading cover art: %v\n", err)
	}

	// Add cover art to m4a file
	// ffmpeg -i spot.flac -i spot-cover.jpeg -map 0 -map 1 -c:a copy -disposition:v attached_pic -metadata:s:v
	//title="Album cover" -metadata:s:v comment="Cover (front)" output.m4a
	log.Printf("Downloading cover art: %s\n", coverArtPath)

	finalSongFilePath := songDirPath + fmt.Sprintf("/%s.m4a", song.Name)

	ffmpegCmd := exec.Command("ffmpeg",
		"-i", songFilePath+".m4a",
		"-i", coverArtPath,
		"-map_metadata", "0",
		"-map", "0",
		"-map", "1",
		"-c", "copy",
		"-disposition:v", "attached_pic",
		finalSongFilePath)

	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr

	log.Printf("Adding album art...")
	log.Printf("%v", ffmpegCmd.Args)
	err = ffmpegCmd.Run()
	if err != nil {

		log.Printf("Error running ffmpeg command: %v\n", err)
		return err
	}

	log.Printf("Audio download, conversion, and metadata embedding completed successfully!")

	return nil

}

func scrapeSongs(db *database.Queries, wg *sync.WaitGroup, song database.Song) {
	defer wg.Done()
	log.Printf("Scraping song: %v", song)

}
