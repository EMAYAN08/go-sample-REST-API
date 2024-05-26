package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items []Item

func main() {
	fmt.Println("Sample implementation of REST APIs")

	items = append(items, Item{ID: "1", Name: "Apple"}, Item{ID: "2", Name: "Bananana"})
	router := mux.NewRouter()
	router.HandleFunc("/items", getAllItems).Methods("GET")
	router.HandleFunc("/item/{id}", getItem).Methods("GET")
	router.HandleFunc("/item", createItem).Methods("POST")
	router.HandleFunc("/item/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")

	fmt.Println("Server started and listening at port: 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}

func getAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range items {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode("Item Not Found!")

}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)
	newItem.ID = uuid.New().String()
	items = append(items, newItem)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, item := range items {
		if item.ID == params["id"] {
			var updateItem Item
			items = append(items[:idx], items[idx+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&updateItem)
			updateItem.ID = params["id"]
			items = append(items, updateItem)
			json.NewEncoder(w).Encode(updateItem)
			return
		}
	}
	json.NewEncoder(w).Encode("Item not found!")

}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, item := range items {
		if item.ID == params["id"] {
			items = append(items[:idx], items[idx+1:]...)
			json.NewEncoder(w).Encode(items)
			break
		}
	}
}
