package backendlogic

import (
	"encoding/json"
	"log"
	"net/http"
)

//----IMPORTANT----
//If two people try to append a book to your slice at
//the exact same nanosecond, your program might crash or "panic."
//-----DEAL WITH IT LATER-----

func GetItemsHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    err := json.NewEncoder(w).Encode(Books)
    if err != nil {
        //In web development, once you start writing to the body(w.Write or Encode),
		//you have technically sent a 200 OK status. If an error happens during the
		//encoding, you can't "take back" the 200 OK.
        log.Printf("Failed to encode books: %v", err)
        return
    }
}

func GetItemHandler(w http.ResponseWriter, r *http.Request){
	IdFromUrl := r.PathValue("id") //route: /books/{id}
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for _, currentBook := range Books {
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
	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    Books = append(Books, newBook)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	//A simple struct with strings and integers is "safe."
	//The only reason json.Encode usually fails is if the data
	//contains things that JSON can't represent
	//(like functions, complex channels, or circular references)
    json.NewEncoder(w).Encode(newBook)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request){
	var updatedBook Book
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
	for i := range Books {
		if Books[i].ID == IdFromUrl {
			updatedBook.ID = IdFromUrl
            Books[i] = updatedBook

			w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Books[i])
            return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request){
	IdFromUrl := r.PathValue("id") //route: /books/{id}
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for i, currentBook := range Books {
		if currentBook.ID == IdFromUrl {
			copy(Books[i:], Books[i+1:])
			//A slice is just a "window" looking at a fixed-size Underlying Array
			//The Garbage Collector (GC) won't clean up that book because it thinks it's still "visible."
			//By setting the last element to Book{} (an empty struct) before truncating, you "zero out" the memory so the GC can reclaim it
			Books[len(Books)-1] = Book{}
            Books = Books[:len(Books)-1]
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

