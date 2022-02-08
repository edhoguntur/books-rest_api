package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

// DB variable to access database
var db *sql.DB

// DB Model
type conn struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Initialized Database
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error Loading .env file: %s \n", err.Error())
		return
	}

	connInfo := conn{
		Host:     os.Getenv("POSTGRES_URL"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	//try open connection with postgresql db using connInfo var
	db, err = sql.Open("postgres", connToString(connInfo))
	if err != nil {
		fmt.Printf("Error connecting to the DB: %s\n", err.Error())
		return
	} else {
		fmt.Printf("PostgreSQL is open\n")
	}

	// check if we can ping our DB
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error could not ping database: %s\n", err.Error())
		return
	} else {
		fmt.Printf("DB pinged successfully and ready to the rock!!!\n")
	}
}

// Take our connection struct and convert to a string for our db connection info
func connToString(info conn) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.Host, info.Port, info.User, info.Password, info.DBName)
}

// Book Model
type Book struct {
	ID     int64   `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Model
type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// GetAllBooks Service
func GetAllBooks() ([]Book, error) {
	var books []Book
	query := `SELECT * FROM BOOKS;`
	rows, err := db.Query(query)
	if err != nil {
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var isbn, title, firstname, lastname string

		err := rows.Scan(&id, &isbn, &title, &firstname, &lastname)
		if err != nil {
			return books, err
		}

		book := Book{ID: id, Isbn: isbn, Title: title, Author: &Author{FirstName: firstname, LastName: lastname}}
		books = append(books, book)
	}
	return books, nil
}

// GetSingleBook Service
func GetSingleBook(id int64) ([]Book, error) {
	var books []Book
	query := `SELECT * FROM BOOKS WHERE ID=$1;`
	rows, err := db.Query(query, id)
	if err != nil {
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var isbn, title, firstname, lastname string

		err := rows.Scan(&id, &isbn, &title, &firstname, &lastname)
		if err != nil {
			return books, err
		}

		book := Book{ID: id, Isbn: isbn, Title: title, Author: &Author{FirstName: firstname, LastName: lastname}}
		books = append(books, book)
	}
	return books, nil
}

// CreateBook Service
func CreateBook(book Book) error {

	query := `INSERT INTO BOOKS (isbn, title, firstname, lastname) values($1, $2, $3, $4);`

	_, err := db.Exec(query, book.Title, book.Isbn, book.Author.FirstName, book.Author.LastName)

	if err != nil {
		return err
	}

	return nil
}

// EditBook Service
func EditBook(id int64, book Book) error {

	query := `UPDATE BOOKS SET isbn=$1, title=$2, firstname=$3, lastname=$4 WHERE id=$5;`

	_, err := db.Exec(query, book.Isbn, book.Title, book.Author.FirstName, book.Author.LastName, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBook Service
func DeleteBook(id int64) error {
	query := `DELETE FROM BOOKS WHERE id=$1;`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

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

	books, err := GetAllBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	} else {
		json.NewEncoder(w).Encode(books)
	}
}

// Get a single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	book, err := GetSingleBook(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(book)
}

// Add a new book
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var book Book
	err := decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = CreateBook(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var book Book
	err := decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	params := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(params, 10, 64)

	err = EditBook(id, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = DeleteBook(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// initialized router
	r := mux.NewRouter()

	// Hardcoded data
	//books = append(books, Book{ID: "1", Isbn: "635483", Title: "Book One", Author: &Author{FirstName: "Guntur", LastName: "Adhitama"}})

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	//Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
