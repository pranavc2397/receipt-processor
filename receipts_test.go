package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessReceipt(t *testing.T) {
	// Create a new HTTP request for the /receipts/process endpoint
	payload := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
			{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
			{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
			{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
			{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"}
		],
		"total": "35.35"
	}`
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(processReceipt)

	// Serve the HTTP request using the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	if _, ok := response["id"]; !ok {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestGetPoints(t *testing.T) {
	// First, process a receipt to get an ID
	payload := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
			{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
			{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
			{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
			{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"}
		],
		"total": "35.35"
	}`
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(processReceipt)
	handler.ServeHTTP(rr, req)

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	receiptID, ok := response["id"]
	if !ok {
		t.Fatalf("handler returned unexpected body: got %v", rr.Body.String())
	}

	// Now, get the points for the processed receipt
	req, err = http.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getPoints)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var pointsResponse map[string]int
	err = json.NewDecoder(rr.Body).Decode(&pointsResponse)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	if _, ok := pointsResponse["points"]; !ok {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	// Verify the points
	expectedPoints := 28
	if pointsResponse["points"] != expectedPoints {
		t.Errorf("handler returned unexpected points: got %v want %v", pointsResponse["points"], expectedPoints)
	}
}

func TestGetPointsInvalidID(t *testing.T) {
	// Attempt to get points for an invalid receipt ID
	req, err := http.NewRequest("GET", "/receipts/invalid-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPoints)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
