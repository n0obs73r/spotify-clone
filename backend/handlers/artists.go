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

type Artist struct {
	Name   string   `json:"name"`
	Albums []string `json:"albums"`
}

func GetArtistsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	filter := query.Get("filter")
	sortBy := query.Get("sort_by")

	sortFuncs := map[string]func(a, b *Artist) bool{
		"name": func(a, b *Artist) bool {
			return a.Name < b.Name
		},
		"albums_count": func(a, b *Artist) bool {
			return len(a.Albums) < len(b.Albums)
		},
	}

	artistsDir := "D:/Music/MP3"

	artistMap := make(map[string]*Artist)

	err := filepath.Walk(artistsDir, func(path string, info os.FileInfo, err error) error {

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

				artist := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

				album := "Unknown Album"

				if _, ok := artistMap[artist]; !ok {

					artistMap[artist] = &Artist{
						Name:   artist,
						Albums: []string{album},
					}
				} else {

					artistMap[artist].Albums = append(artistMap[artist].Albums, album)
				}

				return nil
			}

			artist := mp3File.Artist()
			album := mp3File.Album()

			if _, ok := artistMap[artist]; !ok {

				artistMap[artist] = &Artist{
					Name:   artist,
					Albums: []string{album},
				}
			} else {

				artistMap[artist].Albums = append(artistMap[artist].Albums, album)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking through directory: %v", err)
		http.Error(w, "Error reading artists", http.StatusInternalServerError)
		return
	}

	artists := make([]Artist, 0, len(artistMap))
	for _, artist := range artistMap {
		artists = append(artists, *artist)
	}

	if filter != "" {
		artists = filterArtists(artists, filter)
	}

	if sortBy != "" {
		sortFunc, ok := sortFuncs[sortBy]
		if ok {
			sort.Slice(artists, func(i, j int) bool {
				return sortFunc(&artists[i], &artists[j])
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artists); err != nil {
		log.Printf("Error encoding artists to JSON: %v", err)
		http.Error(w, "Error encoding artists", http.StatusInternalServerError)
		return
	}
}

func filterArtists(artists []Artist, filter string) []Artist {
	filteredArtists := make([]Artist, 0)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(filter)) {

			filteredArtists = append(filteredArtists, artist)
		}
	}

	return filteredArtists
}
