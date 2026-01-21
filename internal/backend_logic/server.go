package backendlogic

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", GetItemsHandler)
	mux.HandleFunc("POST /books", PostItemHandler)

	mux.HandleFunc("GET /books/{id}", GetItemHandler)
	mux.HandleFunc("PUT /books/{id}", UpdateItemHandler)
	mux.HandleFunc("DELETE /books/{id}", DeleteItemHandler)

	return mux
}