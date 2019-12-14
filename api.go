package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type API struct {
	db    *DB
	stats *Stats
}

func (a *API) Run() {
	http.HandleFunc("/add", a.handleAdd)
	http.HandleFunc("/delete", a.handleDelete)
	http.HandleFunc("/list", a.handleList)
	http.HandleFunc("/stats", a.handleStats)
	http.ListenAndServe("localhost:8000", nil)
}

func (a *API) handleAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var addRequest struct {
		URL string `json:"url"`
	}

	err := decoder.Decode(&addRequest)

	if err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return
	}

	if err := a.db.AddURL(addRequest.URL); err != nil {
		log.Printf("Error adding URL: %v", err)
	}
}

func (a *API) handleDelete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var deleteRequest struct {
		URL string `json:"url"`
	}

	err := decoder.Decode(&deleteRequest)

	if err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return
	}

	a.db.DeleteURL(deleteRequest.URL)
}

func (a *API) handleList(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	var listResponse struct {
		URLs []string `json:"urls"`
	}

	listResponse.URLs = a.db.ListURLs()

	err := encoder.Encode(&listResponse)

	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		return
	}
}

func (a *API) handleStats(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	var statsResponse struct {
		Message string `json:"message"`
	}

	statsResponse.Message = a.stats.Format(a.db.Size())

	err := encoder.Encode(&statsResponse)

	if err != nil {
		log.Printf("Failed to encode JSON: %v", err)
		return
	}
}
