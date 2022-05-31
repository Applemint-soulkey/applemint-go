package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("starting server...")

	os.Setenv("env_mongo_id", "rlatmfrl24")
	os.Setenv("env_mongo_pwd", "397love")
	os.Setenv("env_dropbox_access_token", "sl.BIQJDuNkA6ZbRJ5rCfZiKfNb0f3Yx4OBgnLyUDON0-W1mGOmVo18aJWdF2FNy4DokNQS_9QZv1IPwVLBOlkJ0wMq1Zwr66zJRapJbcL_TfW1LqGARMVlUcIOteKdlS0O0J8JsCUy")
	os.Setenv("env_raindrop_access_token", "54b0d37e-03b2-4453-b6d1-b59a49c4b536")

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/item/move/{id}", handleMoveItemRequest).Methods("GET")
	r.HandleFunc("/item/keep/{id}", handleKeepItemRequest).Methods("POST")
	r.HandleFunc("/item/{collection}/{id}", handleItemRequest).Methods("GET", "POST","DELETE")
	r.HandleFunc("/collection/{target}", handleClearCollectionRequest).Methods("DELETE")
	r.HandleFunc("/crawl/{target}", handleCrawlRequest).Methods("GET")
	
	r.HandleFunc("/dropbox/", handleDropboxRequest).Methods("GET")
	r.HandleFunc("/raindrop/{collectionId}", handleRaindropRequest).Methods("POST")
	r.HandleFunc("/raindrop/collections", handleRaindropCollectionRequest).Methods("GET")

	http.Handle("/", r)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if  port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP service.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}