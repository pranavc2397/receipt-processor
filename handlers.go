package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var (
	receipts      = make(map[string]ProcessedReceipt)
	receiptsMutex = &sync.Mutex{}
)

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	points := calculatePoints(receipt)
	id := uuid.New().String()

	receiptsMutex.Lock()
	receipts[id] = ProcessedReceipt{ID: id, Points: points}
	receiptsMutex.Unlock()

	response := map[string]string{"id": id}
	json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/receipts/"):]
	id = id[:len(id)-7]

	receiptsMutex.Lock()
	receipt, exists := receipts[id]
	receiptsMutex.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": receipt.Points}
	json.NewEncoder(w).Encode(response)
}
