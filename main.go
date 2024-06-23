package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var tpl *template.Template

type movie struct {
	ID    int
	Title string
	Year  int
}

var (
	movies = []movie{
		{1, "The Shawshank Redemption", 1994},
		{2, "The Godfather", 1972},
		{3, "The Dark Knight", 2008},
	}
)

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html", "templates/updateForm.html", "templates/addform.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/get", getHandler)
	//http.HandleFunc("/put", putHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tpl.Execute(w, movies)
}

func addHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		yearStr := r.FormValue("year")
		year, err := strconv.Atoi(yearStr)
		var newID int

		if err != nil {
			http.Error(w, "Invalid year", http.StatusBadRequest)
			return
		}

		if len(movies) == 0 {
			newID = 1
		} else {
			newID = movies[len(movies)-1].ID + 1
		}

		newMovie := movie{ID: newID, Title: title, Year: year}

		movies = append(movies, newMovie)

		// Update the list with the new entry
		w.Header().Set("Content-Type", "text/html")
		err = tpl.ExecuteTemplate(w, "list-range", movies)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {

		idStr := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(idStr[len(idStr)-1])

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		for i, m := range movies {
			if m.ID == id {
				movies = append(movies[:i], movies[i+1:]...)
				break
			}
		}
		tpl.ExecuteTemplate(w, "list-range", movies)

	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(idStr[len(idStr)-1])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find the movie with the given ID
	var selectedMovie movie
	for _, movie := range movies {
		if movie.ID == id {
			selectedMovie = movie
			break
		}
	}

	// Render the update form with the movie details
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "updateForm", selectedMovie)
}
