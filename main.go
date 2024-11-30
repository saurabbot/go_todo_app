package main

import (
	"encoding/json" // Add this
	"fmt"

	// Add this
	"net/http"
	"sync"
	"time"
)

type Book struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"publishedYear"`
	Genre         int       `json:"genre"`
	IsAvailable   bool      `json:"isAvailable"`
	AddedAt       time.Time `json:"addedAt"`
}

type BookStore struct {
	books map[string]Book
	mutex sync.RWMutex
}

var store = &BookStore{
	books: make(map[string]Book),
}

func main() {
	http.HandleFunc("/books", handleBooks)

	fmt.Println("Server starting on http://localhost:8080") // Print BEFORE starting server
	err := http.ListenAndServe(":8080", nil)

	// These lines will never be reached unless there's an error
	if err != nil {
		fmt.Println("Server failed to start")
		panic(err)
	}
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w, r)
	case http.MethodPost:
		addBook(w, r) // Changed to match function name
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getBooks(w http.ResponseWriter, _ *http.Request) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	books := make([]Book, 0, len(store.books))
	for _, book := range store.books {
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request) { // Changed to singular
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book.ID = fmt.Sprintf("book_%d", time.Now().UnixNano())
	book.AddedAt = time.Now()
	book.IsAvailable = true

	store.mutex.Lock()
	store.books[book.ID] = book
	store.mutex.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
