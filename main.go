package main

import (
	"applemint-go/crawl"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	crawl.Crawl("http://golang.org/")
	log.Print("starting server...")
	http.HandleFunc("/", handler)

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

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv(("NAME"))
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}