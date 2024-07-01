package main

import (
	"context"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"log"
	"sync"
	"time"
)

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

func scrapeSongs(db *database.Queries, wg *sync.WaitGroup, song database.Song) {
	defer wg.Done()
	log.Printf("Scraping song: %v", song)

}
