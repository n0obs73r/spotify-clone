package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func Search(w http.ResponseWriter, r *http.Request) {

	tracksDir := "D:/Music/MP3"

	query := r.URL.Query().Get("q")
	log.Printf("Search query: %s\n", query)

	searchResults := make([]string, 0)

	err := filepath.Walk(tracksDir, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(strings.ToLower(path), ".mp3") {
			log.Printf("Checking file: %s\n", path)

			file, err := os.Open(path)
			if err != nil {
				log.Printf("Error opening file %s: %v\n", path, err)
				return err
			}
			defer file.Close()

			mp3File, err := tag.ReadFrom(file)
			if err != nil {
				if err == tag.ErrNoTagsFound {
					log.Printf("No tags found in file %s, using filename as title\n", path)

					title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
					searchResults = append(searchResults, title)
					return nil
				}
				log.Printf("Error parsing file %s: %v\n", path, err)
				return err
			}

			title := mp3File.Title()
			artist := mp3File.Artist()
			album := mp3File.Album()

			log.Printf("Metadata - Title: %s, Artist: %s, Album: %s\n", title, artist, album)

			if strings.Contains(strings.ToLower(title), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(artist), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(album), strings.ToLower(query)) {
				log.Printf("Match found for query '%s' in file: %s\n", query, path)
				searchResults = append(searchResults, title)
			} else {

				filename := filepath.Base(path)
				if strings.Contains(strings.ToLower(filename), strings.ToLower(query)) {
					log.Printf("Match found for query '%s' in filename: %s\n", query, filename)
					searchResults = append(searchResults, filename)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error searching tracks: %v\n", err)
		http.Error(w, "Error searching tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}
