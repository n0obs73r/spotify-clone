package main

import (
	"log"
	"net/http"

	"backend/handlers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	r := mux.NewRouter()

	// Define CORS middleware function
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}

	// Apply CORS middleware to all routes
	r.Use(corsMiddleware)

	// Define your routes
	r.HandleFunc("/playlists", handlers.GetPlaylists).Methods("GET")
	r.HandleFunc("/albums", handlers.GetAlbumsHandler).Methods("GET")
	r.HandleFunc("/artists", handlers.GetArtistsHandler).Methods("GET")
	r.HandleFunc("/genres", handlers.GetGenresHandler).Methods("GET")
	r.HandleFunc("/search", handlers.Search).Methods("GET")
	r.HandleFunc("/tracks", handlers.GetTracksHandler).Methods("GET")

	r.HandleFunc("/albums/{title}/art", handlers.GetAlbumArt).Methods("GET")
	r.HandleFunc("/albums/{title}/songs", handlers.GetAlbumSongsHandler).Methods("GET")

	r.PathPrefix("/songs/").Handler(http.StripPrefix("/songs/", http.FileServer(http.Dir("./mp3"))))

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", n)
}
