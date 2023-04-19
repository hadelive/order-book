package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	lib "github.com/hadelive/order-book/lib"
)

func main() {
	orderBook := lib.NewOrderBook()

	// Define the HTTP endpoints and handlers
	http.HandleFunc("/addOrder", func(w http.ResponseWriter, r *http.Request) {
		priceStr := r.URL.Query().Get("price")
		quantityStr := r.URL.Query().Get("quantity")

		price, err := strconv.Atoi(priceStr)
		if err != nil {
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		var orderID string
		side := r.URL.Query().Get("side")
		if side == "buy" {
			orderID = orderBook.BuyHeap.AddOrder(price, quantity)
		} else if side == "sell" {
			orderID = orderBook.SellHeap.AddOrder(price, quantity)
		} else {
			http.Error(w, "Invalid side", http.StatusBadRequest)
			return
		}

		response := struct {
			OrderID string `json:"order_id"`
		}{OrderID: orderID}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonResponse)
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
