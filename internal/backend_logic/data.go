package backendlogic

import (
	"sync"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

type BookStore struct {
    mux sync.RWMutex
    Books []models.Book
}
var store = BookStore{
    Books: []models.Book{},
}
//Go has a special function called init(). It runs automatically
//once when the package is first loaded, before the main() function even starts.
//This is useful if you want to generate IDs or do more complex setup.
func init() {
    store.Books = append(store.Books, models.Book{ID: "1", Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015})
}