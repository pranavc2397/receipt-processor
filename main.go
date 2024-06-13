package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/receipts/process", processReceipt).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

    log.Println("Server starting on port 3000")
    if err := http.ListenAndServe(":3000", r); err != nil {
        log.Fatal(err)
    }
}
