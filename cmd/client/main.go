package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

func main(){
	b1 := models.Book{ID: "1", Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015}
	b2 := models.Book{ID: "2", Title: "Learning Go", Author: "Miek Gieben", Year: 2019}
	postRequest(b1)
	postRequest(b2)
	getRequest()
}


func postRequest(s models.Book){
	jsonData, err := json.Marshal(s)
	if err != nil{ log.Fatal("Err")}
	_, er := http.Post("http://localhost:8080/books", "application/json", bytes.NewReader(jsonData))
	if er!= nil { log.Fatal("Err")}
	fmt.Println("Book added successfully!")
}

func getRequest(){
	resp, err := http.Get("http://localhost:8080/books")
	if err!=nil{ log.Fatal("Err")}
	data ,err := io.ReadAll(resp.Body)
	var books []models.Book
	err = json.Unmarshal(data, &books)
	if err!= nil {log.Fatal("Err")}
	fmt.Println("Books :")
	for i := range books{
		fmt.Println("Title:", books[i].Title," Author:", books[i].Author, " Year: ", books[i].Year)
	}

}