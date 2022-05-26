package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("starting server...")

	// os.Setenv("env_mongo_id", "rlatmfrl24")
	// os.Setenv("env_mongo_pwd", "397love")

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/item/move/{id}", handleMoveItemRequest).Methods("GET")
	r.HandleFunc("/item/keep/{id}", handleKeepItemRequest).Methods("POST")
	r.HandleFunc("/item/{collection}/{id}", handleItemRequest).Methods("GET", "POST","DELETE")
	r.HandleFunc("/collection/{target}", handleClearCollectionRequest).Methods("DELETE")
	r.HandleFunc("/crawl/{target}", handleCrawlRequest).Methods("GET")
	
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
