package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Book Model
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Model
type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Init Book variable as slice book struct
var books []Book

// test server
func index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(`Hello World`)
	if err != nil {
		return
	}
}

// Get all Books
func getBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		return
	}
}

// Get a single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// loop through books and find requested book
	for _, book := range books {
		if book.ID == params["id"] {
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				return
			}
		}
	}
	err := json.NewEncoder(w).Encode(&Book{})
	if err != nil {
		return
	}
}

// Add a new book
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000)) // Mock ID  will change after added DB
	books = append(books, book)
	err := json.NewEncoder(w).Encode(book)
	if err != nil {
		return
	}
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[:index+1]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				return
			}
			return
		}
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[:index+1]...)
			break
		}
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		return
	}
}

func main() {
	// initialized router
	r := mux.NewRouter()

	// Hardcoded data
	books = append(books, Book{ID: "1", Isbn: "635483", Title: "Book One", Author: &Author{FirstName: "Guntur", LastName: "Adhitama"}})

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id", deleteBook).Methods("DELETE")

	//Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
