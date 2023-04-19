package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Rated string `json:"rated"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Title: "The Shawshank Redemption", Rated: "R"})
	movies = append(movies, Movie{ID: "2", Title: "The Godfather", Rated: "R"})
	movies = append(movies, Movie{ID: "3", Title: "The Dark Knight", Rated: "R"})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	json.NewEncoder(w).Encode(&Movie{})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	if movie.ID == "" || movie.Title == "" {
		http.Error(w, "Empty ID or title not allowed", http.StatusBadRequest)
		return
	}

	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for ind, movie := range movies {
		if movie.ID == params["id"] {
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)

			if movie.ID == "" || movie.Title == "" {
				http.Error(w, "Empty ID or title not allowed", http.StatusBadRequest)
				return
			}

			movies = append(movies[:ind], movies[ind+1:]...)

			params["id"] = movie.ID
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for ind, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:ind], movies[ind+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}
