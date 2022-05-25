package main

import (
	"applemint-go/crawl"
	"applemint-go/crud"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("starting server...")
	os.Setenv("env_mongo_id", "rlatmfrl24")
	os.Setenv("env_mongo_pwd", "397love")
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/item/move/", handleMoveItemRequest).Methods("POST")
	r.HandleFunc("/item/{collection}/{id}", handleDeleteItemRequest).Methods("DELETE")
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

func handleCrawlRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	switch expression := mux.Vars(r)["target"]; expression {
	case "bp":
		json.NewEncoder(w).Encode(crawl.CrawlBP())
	case "isg":
		json.NewEncoder(w).Encode(crawl.CrawlISG())
	}
}

func handleMoveItemRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func handleDeleteItemRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	targetId := mux.Vars(r)["id"]
	targetCollection := mux.Vars(r)["collection"]
	delCnt := crud.DeleteItem(targetId, targetCollection)
	if delCnt > 0 {
		fmt.Fprintf(w, "{\"msg\": \"item deleted from %s -> %s\"}", targetCollection, targetId)
	} else {
		fmt.Fprintf(w, "{\"error\": \"cannot find item from %s -> %s\"}", targetCollection, targetId)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv(("NAME"))
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}