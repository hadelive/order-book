package main

// func main() {
// 	orderBook := lib.NewOrderBook()

// 	// Define the HTTP endpoints and handlers
// 	http.HandleFunc("/addOrder", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var req struct {
// 			Price    int    `json:"Price"`
// 			Quantity int    `json:"Quantity"`
// 			Side     string `json:"Side"`
// 		}

// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, "Invalid request", http.StatusBadRequest)
// 			return
// 		}

// 		defer r.Body.Close()
// 		var OrderID string
// 		if req.Side == "buy" {
// 			OrderID = orderBook.BuyHeap.AddOrder(req.Price, req.Quantity)
// 		} else if req.Side == "sell" {
// 			OrderID = orderBook.SellHeap.AddOrder(req.Price, req.Quantity)
// 		} else {
// 			http.Error(w, "Invalid Side", http.StatusBadRequest)
// 			return
// 		}

// 		response := struct {
// 			OrderID string `json:"order_id"`
// 		}{OrderID: OrderID}

// 		jsonResponse, err := json.Marshal(response)
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		orderBook.MatchOrders()

// 		w.Header().Set("Content-type", "application/json")
// 		w.Write(jsonResponse)
// 	})

// 	http.HandleFunc("/cancelOrder", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodDelete {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var req struct {
// 			OrderID string `json:"order_id"`
// 			Side    string `json:"Side"`
// 		}

// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			log.Printf("Error decoding request body: %s", err)
// 			http.Error(w, "Invalid request", http.StatusBadRequest)
// 			return
// 		}

// 		if req.Side == "buy" {
// 			orderBook.BuyHeap.CancelOrder(req.OrderID)
// 		} else if req.Side == "sell" {
// 			orderBook.SellHeap.CancelOrder(req.OrderID)
// 		} else {
// 			http.Error(w, "Invalid Side", http.StatusBadRequest)
// 			return
// 		}

// 		fmt.Fprintf(w, "Order %s canceled successfully", req.OrderID)
// 	})

// 	// Start the HTTP server
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
