package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

const baseURL = "http://localhost:8080/books"

// httpClient is defined at the package level to enable HTTP Connection Pooling.
// Reusing a single client avoids the expensive overhead of creating new TCP
// connections (handshakes) for every request, significantly boosting performance.
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}


func main(){
	b1 := models.Book{ID: "2", Title: "Learning Go", Author: "Miek Gieben", Year: 2019}
	postRequest(b1)
	getAllRequest()
	updated := models.Book{Title: "Learning Go", Author: "John Doe", Year: 2019}
	updateBook("2", updated)
	getBook("2")
	deleteBook("1")
	getAllRequest()
}


func postRequest(s models.Book){
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(s)
	if err != nil {
		log.Printf("Error encoding book: %v", err)
		return
	}
	resp, err := httpClient.Post(baseURL, "application/json", &buf)
	if err != nil {
		log.Printf("Network error: %v", err)
		return
	}
	// Draining ensures any remaining response bytes are read so the
    // connection is "clean" and can be returned to the pool for reuse.
	defer func() {
   		io.Copy(io.Discard, resp.Body)
    	resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Book added successfully!")
	} else {
		fmt.Printf("Server returned error: %s\n", resp.Status)
	}
}

func getAllRequest(){
	resp, err := httpClient.Get(baseURL)
	if err!=nil{
		log.Printf("Network Error: %v", err)
        return
	}
	defer func() {
    	io.Copy(io.Discard, resp.Body)
    	resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
        fmt.Printf("Server returned error: %s\n", resp.Status)
        return
    }
	var books []models.Book
	err = json.NewDecoder(resp.Body).Decode(&books)
    if err != nil {
        log.Printf("Decoding Error: %v", err)
        return
    }
	fmt.Println("\n--- Current Book List ---")
	for _, b := range books {
        fmt.Printf("Title: %-20s | Author: %-15s | Year: %d\n", b.Title, b.Author, b.Year)
    }

}

func deleteBook(id string) {
	// http.Get and http.Post are shortcuts.
	// For DELETE and PUT, we must create a "Request" object manually.
	url := fmt.Sprintf("%s/%s", baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Send the request using the default client
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending delete: %v", err)
		return
	}
	defer func() {
   		io.Copy(io.Discard, resp.Body)
    	resp.Body.Close()
	}()

	switch resp.StatusCode {
    case http.StatusNoContent, http.StatusOK:
        fmt.Printf("Successfully deleted book [%s]\n", id)
    case http.StatusNotFound:
        fmt.Printf("Delete failed: Book [%s] not found\n", id)
    default:
        fmt.Printf("Unexpected error: %s\n", resp.Status)
    }
}

func updateBook(id string, updatedBook models.Book) {
	url := fmt.Sprintf("%s/%s", baseURL, id)
	var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(updatedBook); err != nil {
        log.Printf("Encoding error: %v", err)
        return
    }
	req, err := http.NewRequest(http.MethodPut, url, &buf)
	if err != nil {
		log.Printf("Error creating update request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending update: %v", err)
		return
	}
	defer func() {
   		io.Copy(io.Discard, resp.Body) 
    	resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successfully updated book [%s]\n", id)
	} else {
		fmt.Printf("Update failed with status: %s\n", resp.Status)
	}
}

func getBook(id string) {
	url := fmt.Sprintf("%s/%s", baseURL, id)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("Error fetching book: %v", err)
		return
	}
	defer func() {
   		io.Copy(io.Discard, resp.Body) 
    	resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
        if resp.StatusCode == http.StatusNotFound {
            fmt.Printf("Book [%s] not found.\n", id)
        } else {
            fmt.Printf("Server error: %s\n", resp.Status)
        }
        return
    }

	var book models.Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		log.Printf("Error decoding: %v", err)
		return
	}

	fmt.Printf("\n--- Details for Book [%s] ---\n", id)
    fmt.Printf("%-10s %s\n", "Title:", book.Title)
    fmt.Printf("%-10s %s\n", "Author:", book.Author)
    fmt.Printf("%-10s %d\n", "Year:", book.Year)
}