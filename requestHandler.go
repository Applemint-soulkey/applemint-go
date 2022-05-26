package main

import (
	"applemint-go/crawl"
	"applemint-go/crud"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleClearCollectionRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	target := mux.Vars(r)["target"]
	delCnt := crud.ClearCollection(target)
	if delCnt > 0 {
		fmt.Fprintf(w, "Deleted %d items from collection %s", delCnt, target)
	} else {
		fmt.Fprintf(w, "Collection %s is empty", target)
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
	targetId := mux.Vars(r)["id"]
	target_coll := r.URL.Query().Get("target")
	origin_coll := r.URL.Query().Get("origin")
	if target_coll == "" || origin_coll == "" {
		fmt.Fprintf(w, "Missing parameters")
		return
	}
	err := crud.MoveItem(targetId, origin_coll, target_coll)
	if err != nil {
		fmt.Fprintf(w, "Error moving item: %s", err)
		return
	}
	fmt.Fprintf(w, "Item moved from %s to %s", origin_coll, target_coll)
}

func handleKeepItemRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	targetId := mux.Vars(r)["id"]
	item := crud.Item{}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		fmt.Fprintf(w, "Error decoding item: %s", err)
		return
	}
	updateCnt := crud.UpdateItem(targetId, "new", item)
	if updateCnt > 0 {
		fmt.Fprintf(w, "Updated %d items\n", updateCnt)
	} else {
		fmt.Fprintf(w, "No items updated")
		return
	}
	err = crud.MoveItem(targetId, "new", "keep")
	if err != nil {
		fmt.Fprintf(w, "Error moving item: %s", err)
		return
	}

	fmt.Fprintf(w, "Updated Item moved from new to keep")
}

func handleItemRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	targetId := mux.Vars(r)["id"]
	targetCollection := mux.Vars(r)["collection"]
	switch r.Method {
	case "GET":
		item, err := crud.GetItem(targetId, targetCollection)
		if err != nil {
			fmt.Fprintf(w, "Error getting item: %s", err)
			return
		}
		json.NewEncoder(w).Encode(item)

	case "POST":
		item := crud.Item{}
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			fmt.Fprintf(w, "Error decoding item: %s", err)
			return
		}
		updateCnt := crud.UpdateItem(targetId, targetCollection, item)
		if updateCnt > 0 {
			fmt.Fprintf(w, "Updated %d items from collection %s", updateCnt, targetCollection)
		} else {
			fmt.Fprintf(w, "Collection %s is empty", targetCollection)
		}
	case "DELETE":
		delCnt := crud.DeleteItem(targetId, targetCollection)
		if delCnt > 0 {
			fmt.Fprintf(w, "{\"msg\": \"item deleted from %s -> %s\"}", targetCollection, targetId)
		} else {
			fmt.Fprintf(w, "{\"error\": \"cannot find item from %s -> %s\"}", targetCollection, targetId)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv(("NAME"))
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}