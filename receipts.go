package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type ProcessedReceipt struct {
	ID     string
	Points int
}

// Utility function to calculate points for a receipt
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points += len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if total == math.Floor(total) {
			points += 50
		}
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	if purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if purchaseDate.Day()%2 != 0 {
			points += 6
		}
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if purchaseTime.Hour() == 14 {
			points += 10
		}
	}

	return points
}
