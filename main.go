package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
)

type GroceryItem struct {
	Id       string
	Name     string
	Quantity int
}

var groceryList []GroceryItem = []GroceryItem{}
var nextId = 0

func findItem(g []GroceryItem, id string) *GroceryItem {
	index := slices.IndexFunc(g, func(item GroceryItem) bool {
		return item.Id == id
	})
	if index == -1 {
		return nil
	}
	return &g[index]
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if err != nil {
		body, _ := json.Marshal(map[string]interface{}{
			"error": err.Error(),
		})
		w.Write(body)
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, payload any) {
	w.WriteHeader(statusCode)
	body, err := json.Marshal(payload)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(body)
}

func GetAllItems(w http.ResponseWriter, req *http.Request) {
	writeResponse(w, http.StatusOK, groceryList)
}

func AddItem(w http.ResponseWriter, req *http.Request) {
	var params GroceryItem
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	nextId++
	params.Id = fmt.Sprintf("item-%d", nextId)
	groceryList = append(groceryList, params)
	writeResponse(w, http.StatusOK, params)
}

func GetItem(w http.ResponseWriter, req *http.Request) {
	item := findItem(groceryList, req.PathValue("id"))
	if item == nil {
		writeErrorResponse(w, http.StatusNotFound, nil)
		return
	}
	writeResponse(w, http.StatusOK, *item)
}

func UpdateItem(w http.ResponseWriter, req *http.Request) {
	item := findItem(groceryList, req.PathValue("id"))
	if item == nil {
		writeErrorResponse(w, http.StatusNotFound, nil)
		return
	}

	var params GroceryItem
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if params.Name != "" {
		item.Name = params.Name
	}
	if params.Quantity != 0 {
		item.Quantity = params.Quantity
	}
	writeResponse(w, http.StatusOK, item)
}

func DeleteItem(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	groceryList = slices.DeleteFunc(groceryList, func(item GroceryItem) bool {
		return item.Id == id
	})
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", GetAllItems)
	mux.HandleFunc("POST /items", AddItem)
	mux.HandleFunc("GET /items/{id}", GetItem)
	mux.HandleFunc("PUT /items/{id}", UpdateItem)
	mux.HandleFunc("DELETE /items/{id}", DeleteItem)

	server := &http.Server{
		Addr:    ":8008",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
