package backendlogic

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request){
	time.Sleep(3 * time.Second)
	store.mux.RLock()
	defer store.mux.RUnlock()
    w.Header().Set("Content-Type", "application/json")
    err := json.NewEncoder(w).Encode(store.Books)
    if err != nil {
        //In web development, once you start writing to the body(w.Write or Encode),
		//you have technically sent a 200 OK status. If an error happens during the
		//encoding, you can't "take back" the 200 OK.
        log.Printf("Failed to encode books: %v", err)
        return
    }
	//test proxy headers
	log.Printf("Incoming Request | Host: %s | Real IP: %s",
           r.Header.Get("X-Forwarded-Host"),
           r.Header.Get("X-Forwarded-For"))
}

func GetItemHandler(w http.ResponseWriter, r *http.Request){
	store.mux.RLock()
	defer store.mux.RUnlock()
	IdFromUrl := r.PathValue("id") //route: /books/{id}
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for _, currentBook := range store.Books {
		if currentBook.ID == IdFromUrl {
			w.Header().Set("Content-Type", "application/json")
        	err := json.NewEncoder(w).Encode(currentBook)
			if err != nil {
				log.Printf("Failed to encode books: %v", err)
				return
    		}
        	return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)

}

func PostItemHandler(w http.ResponseWriter, r *http.Request){
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
	store.mux.Lock()
	defer store.mux.Unlock()
    store.Books = append(store.Books, newBook)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	//A simple struct with strings and integers is "safe."
	//The only reason json.Encode usually fails is if the data
	//contains things that JSON can't represent
	//(like functions, complex channels, or circular references)
    json.NewEncoder(w).Encode(newBook)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request){
	var updatedBook models.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
	IdFromUrl := r.PathValue("id") //route: /books/{id}
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	store.mux.Lock()
	defer store.mux.Unlock()
	for i := range store.Books {
		if store.Books[i].ID == IdFromUrl {
			updatedBook.ID = IdFromUrl
            store.Books[i] = updatedBook

			w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(store.Books[i])
            return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request){
	store.mux.Lock()
	defer store.mux.Unlock()
	IdFromUrl := r.PathValue("id") //route: /books/{id}
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for i, currentBook := range store.Books {
		if currentBook.ID == IdFromUrl {
			copy(store.Books[i:], store.Books[i+1:])
			//A slice is just a "window" looking at a fixed-size Underlying Array
			//The Garbage Collector (GC) won't clean up that book because it thinks it's still "visible."
			//By setting the last element to models.Book{} (an empty struct) before truncating, you "zero out" the memory so the GC can reclaim it
			store.Books[len(store.Books)-1] = models.Book{}
            store.Books = store.Books[:len(store.Books)-1]
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

