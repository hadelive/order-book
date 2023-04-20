package main

import (
	"encoding/json"
	"log"
	"net/http"

	lib "github.com/hadelive/order-book/lib"
)

func main() {
	orderBook := lib.NewOrderBook()

	// Define the HTTP endpoints and handlers
	http.HandleFunc("/addOrder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Price    int    `json:"price"`
			Quantity int    `json:"quantity"`
			Side     string `json:"side"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		var orderID string
		if req.Side == "buy" {
			orderID = orderBook.BuyHeap.AddOrder(req.Price, req.Quantity)
		} else if req.Side == "sell" {
			orderID = orderBook.SellHeap.AddOrder(req.Price, req.Quantity)
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

		orderBook.MatchOrders()

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonResponse)
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
