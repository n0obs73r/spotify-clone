package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
)

type Genre struct {
	Name   string   `json:"name"`
	Tracks []string `json:"tracks"`
}

func GetGenresHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("filter")
	sortBy := query.Get("sort_by")

	genresDir := "D:/Music/MP3"
	genreMap := make(map[string]*Genre)

	err := filepath.Walk(genresDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(path), ".mp3") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			mp3File, err := tag.ReadFrom(file)
			if err != nil {
				if err != tag.ErrNoTagsFound {
					log.Printf("Error reading metadata from file %s: %v", path, err)
				}

				genre := "Unknown Genre"

				if _, ok := genreMap[genre]; !ok {
					genreMap[genre] = &Genre{
						Name:   genre,
						Tracks: []string{},
					}
				}
				genreMap[genre].Tracks = append(genreMap[genre].Tracks, filepath.Base(path))

				return nil
			}

			genre := mp3File.Genre()

			if _, ok := genreMap[genre]; !ok {
				genreMap[genre] = &Genre{
					Name:   genre,
					Tracks: []string{},
				}
			}
			genreMap[genre].Tracks = append(genreMap[genre].Tracks, mp3File.Title())
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking through directory: %v", err)
		http.Error(w, "Error reading genres", http.StatusInternalServerError)
		return
	}

	genres := make([]Genre, 0, len(genreMap))
	for _, genre := range genreMap {
		genres = append(genres, *genre)
	}

	if filter != "" {
		genres = filterGenres(genres, filter)
	}

	if sortBy != "" {
		sortGenres(genres, sortBy)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(genres); err != nil {
		log.Printf("Error encoding genres to JSON: %v", err)
		http.Error(w, "Error encoding genres", http.StatusInternalServerError)
		return
	}
}

func filterGenres(genres []Genre, filter string) []Genre {
	filteredGenres := make([]Genre, 0)

	for _, genre := range genres {
		if strings.Contains(strings.ToLower(genre.Name), strings.ToLower(filter)) {
			filteredGenres = append(filteredGenres, genre)
		}
	}

	return filteredGenres
}

func sortGenres(genres []Genre, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(genres, func(i, j int) bool {
			return genres[i].Name < genres[j].Name
		})
	}
}
