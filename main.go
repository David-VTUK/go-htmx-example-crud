package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var tpl *template.Template

type movie struct {
	Title string
	Year  int
}

var (
	movies = []movie{
		{"The Shawshank Redemption", 1994},
		{"The Godfather", 1972},
		{"The Dark Knight", 2008},
	}
)

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// ExecuteTemplate writes the template to the response writer
	//tpl.ExecuteTemplate(w, "index.html", movies)
	tpl.Execute(w, movies)
}

func addHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		yearStr := r.FormValue("year")
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year", http.StatusBadRequest)
			return
		}

		fmt.Printf("Title: %s, Year: %d", title, year)

		newMovie := movie{Title: title, Year: year}

		movies = append(movies, newMovie)
		tpl.ExecuteTemplate(w, "list-range", movies)

	}
}
