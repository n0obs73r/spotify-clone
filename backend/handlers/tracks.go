package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
)

type Track struct {
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Album  string   `json:"album"`
	Genres []string `json:"genres"`
}

func GetTracksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("filter")
	sortBy := query.Get("sort_by")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	tracksDir := "D:/Music/MP3"

	tracks := make([]Track, 0)

	err := filepath.Walk(tracksDir, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(strings.ToLower(path), ".mp3") {

			file, err := os.Open(path)
			if err != nil {
				log.Printf("Error opening file %s: %v", path, err)
				return nil
			}
			defer file.Close()

			mp3File, err := tag.ReadFrom(file)
			if err != nil {
				if err != tag.ErrNoTagsFound {
					log.Printf("Error reading metadata from file %s: %v", path, err)
				}

				title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

				track := Track{
					Title: title,
				}

				tracks = append(tracks, track)

				return nil
			}

			title := mp3File.Title()
			artist := mp3File.Artist()
			album := mp3File.Album()
			genre := mp3File.Genre()

			track := Track{
				Title:  title,
				Artist: artist,
				Album:  album,
				Genres: []string{genre},
			}

			tracks = append(tracks, track)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking through directory: %v", err)
		http.Error(w, "Error reading tracks", http.StatusInternalServerError)
		return
	}

	if filter != "" {
		tracks = filterTracks(tracks, filter)
	}

	if sortBy != "" {
		sortTracks(tracks, sortBy)
	}

	if pageStr != "" && limitStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Printf("Error converting page number to integer: %v", err)
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Error converting limit number to integer: %v", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		tracks = paginateTracks(tracks, page, limit)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tracks); err != nil {
		log.Printf("Error encoding tracks to JSON: %v", err)
		http.Error(w, "Error encoding tracks", http.StatusInternalServerError)
		return
	}
}

func filterTracks(tracks []Track, filter string) []Track {
	filteredTracks := make([]Track, 0)

	for _, track := range tracks {
		if strings.Contains(strings.ToLower(track.Title), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(track.Artist), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(track.Album), strings.ToLower(filter)) ||
			containsGenre(track.Genres, filter) {
			filteredTracks = append(filteredTracks, track)
		}
	}

	return filteredTracks
}

func containsGenre(genres []string, filter string) bool {
	for _, genre := range genres {
		if strings.Contains(strings.ToLower(genre), strings.ToLower(filter)) {
			return true
		}
	}
	return false
}

func sortTracks(tracks []Track, sortBy string) {
	switch sortBy {
	case "title":

		sort.Slice(tracks, func(i, j int) bool {
			return tracks[i].Title < tracks[j].Title
		})
	case "artist":

		sort.Slice(tracks, func(i, j int) bool {
			return tracks[i].Artist < tracks[j].Artist
		})
	case "album":

		sort.Slice(tracks, func(i, j int) bool {
			return tracks[i].Album < tracks[j].Album
		})
	}
}

func paginateTracks(tracks []Track, page, limit int) []Track {
	start := (page - 1) * limit
	end := start + limit

	if start >= len(tracks) {
		return []Track{}
	}
	if end > len(tracks) {
		end = len(tracks)
	}

	return tracks[start:end]
}
