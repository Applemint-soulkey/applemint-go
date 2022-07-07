package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Print("starting server...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Println("Error loading .env file")
	}

	r := mux.NewRouter()

	// CRUD API for item
	r.HandleFunc("/item/move/{id}", handleMoveItemRequest).Methods("GET")
	r.HandleFunc("/item/keep/{id}", handleKeepItemRequest).Methods("POST")
	r.HandleFunc("/item/{collection}/{id}", handleItemRequest).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/items/{collection}", handleItemsRequest).Methods("GET")
	r.HandleFunc("/item/clean", handleClearOldItemsRequest).Methods("GET")

	// Collection Info API
	r.HandleFunc("/collection", handleCollectionRequest).Methods("GET")
	r.HandleFunc("/collection/{target}", handleClearCollectionRequest).Methods("DELETE")
	r.HandleFunc("/collection/info/{collection}", handleCollectionInfoRequest).Methods("GET")

	// Crawl API
	r.HandleFunc("/crawl/{target}", handleCrawlRequest).Methods("GET")

	// Bookmark API
	r.HandleFunc("/item/bookmark", handleBookmarkRequest).Methods("GET", "POST")

	// Gallery API
	r.HandleFunc("/gallery/imgur", handleImgurAnalyzeRequest).Methods("GET")
	r.HandleFunc("/gallery", handleGalleryRequest).Methods("GET")

	// External App API
	r.HandleFunc("/dropbox", handleDropboxRequest).Methods("GET")
	r.HandleFunc("/raindrop/{collectionId}", handleRaindropRequest).Methods("POST")
	r.HandleFunc("/raindrop/collections", handleRaindropCollectionRequest).Methods("GET")

	http.Handle("/", r)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "origin", "x-csrftoken", "x-access-token"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Start HTTP service.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(originsOK, headersOK, methodsOK)(r)); err != nil {
		log.Fatal(err)
	}
}
