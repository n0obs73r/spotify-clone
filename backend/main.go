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

	r.HandleFunc("/playlists", handlers.GetPlaylists).Methods("GET")
	r.HandleFunc("/albums", handlers.GetAlbumsHandler).Methods("GET")
	r.HandleFunc("/artists", handlers.GetArtistsHandler).Methods("GET")
	r.HandleFunc("/genres", handlers.GetGenresHandler).Methods("GET")
	r.HandleFunc("/search", handlers.Search).Methods("GET")
	r.HandleFunc("/tracks", handlers.GetTracksHandler).Methods("GET")

	r.HandleFunc("/albums/{title}/art", handlers.GetAlbumArt).Methods("GET")

	r.PathPrefix("/songs/").Handler(http.StripPrefix("/songs/", http.FileServer(http.Dir("./mp3"))))

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", n)
}
