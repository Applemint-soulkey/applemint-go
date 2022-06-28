package main

import (
	"applemint-go/crawl"
	"encoding/json"
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

	r.HandleFunc("/", handler)
	r.HandleFunc("/item/move/{id}", handleMoveItemRequest).Methods("GET")
	r.HandleFunc("/item/keep/{id}", handleKeepItemRequest).Methods("POST")
	r.HandleFunc("/item/{collection}/{id}", handleItemRequest).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/items/{collection}", handleItemsRequest).Methods("GET")

	r.HandleFunc("/collection/{target}", handleClearCollectionRequest).Methods("DELETE")
	r.HandleFunc("/collection/info/{collection}", handleCollectionInfoRequest).Methods("GET")
	r.HandleFunc("/crawl/{target}", handleCrawlRequest).Methods("GET")

	r.HandleFunc("/item/bookmark", handleBookmarkRequest).Methods("GET", "POST")

	r.HandleFunc("/gallery/imgur", handleImgurAnalyzeRequest).Methods("GET")

	r.HandleFunc("/dropbox/", handleDropboxRequest).Methods("GET")
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

func handleImgurAnalyzeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handleImgurAnalyzeRequest:", r.URL.Path)
	imgurLink := r.URL.Query().Get("link")
	if imgurLink == "" {
		log.Println("handleImgurAnalyzeRequest: missing link")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	images, err := crawl.HandleImgurLink(imgurLink)
	if err != nil {
		log.Println("handleImgurAnalyzeRequest:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("handleImgurAnalyzeRequest:", images)
	w.WriteHeader(http.StatusOK)
	json, err := json.Marshal(images)
	if err != nil {
		log.Println("handleImgurAnalyzeRequest:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}
