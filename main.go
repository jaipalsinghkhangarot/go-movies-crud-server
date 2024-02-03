package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

var movies []Movie

func main() {
	route := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "43782", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	route.HandleFunc("/movies", getAllMovies).Methods("GET")
	route.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	route.HandleFunc("/movies", createNewMovie).Methods("POST")
	route.HandleFunc("/movies/{id}", updatemovie).Methods("PUT")
	route.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("STARTING SERVER AT PORT 8000!!")
	log.Fatal(http.ListenAndServe(":8000", route))

}

func getAllMovies(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "method is not supported", http.StatusNotFound)
	}
	if request.URL.Path != "/movies" {
		http.Error(response, "404 Not Found", http.StatusNotFound)
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(movies)
}

func getMovieById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	if request.Method != "GET" {
		http.Error(response, "method not supported", http.StatusNotFound)
	}

	for index, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(response).Encode(movies[index])
			return
		}
		// for _,item := range movies {
		// 	if item.ID == params["id"]{
		// 		json.NewEncoder(response).Encode(item)
		// 		return
		// 	}
		// }
	}
	http.Error(response, "404 Not Found", http.StatusNotFound)
}

func createNewMovie(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "method not supported", http.StatusNotFound)
	}
	response.Header().Set("Content-Type", "application/json")
	var newMovie Movie
	_ = json.NewDecoder(request.Body).Decode(&newMovie)
	newMovie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, newMovie)
	json.NewEncoder(response).Encode(newMovie)
}

func updatemovie(response http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		http.Error(response, "method not supported", http.StatusNotFound)
	}
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var updateMovie Movie
			_ = json.NewDecoder(request.Body).Decode(&updateMovie)
			updateMovie.ID = params["id"]
			movies = append(movies, updateMovie)
			json.NewEncoder(response).Encode(updateMovie)
			return
		}
	}
	http.Error(response, "404 Not Found", http.StatusNotFound)
}

func deleteMovie(response http.ResponseWriter, request *http.Request) {
	if request.Method != "DELETE" {
		http.Error(response, "method not supported", http.StatusNotFound)
	}

	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(response).Encode(movies)
			return
		}
	}
	http.Error(response, "404 Not Found", http.StatusNotFound)
}
