package orderbook

import (
	"fmt"
	"testing"
)

func TestOrderBook(t *testing.T) {
	// create a new order book
	orderBook := OrderBook{
		buyHeap: OrderHeap{
			less: func(o1, o2 Order) bool { return o1.price > o2.price }, // use greater than for buy orders
		},
		sellHeap: OrderHeap{
			less: func(o1, o2 Order) bool { return o1.price < o2.price }, // use less than for sell orders
		},
	}

	// add buy and sell orders to the order book
	buyOrderID := orderBook.buyHeap.AddOrder(10, 100)
	orderBook.buyHeap.AddOrder(12, 100)
	orderBook.buyHeap.AddOrder(11, 120)
	orderBook.buyHeap.AddOrder(13, 120)

	// Add sell orders to the order book
	sellOrderID1 := orderBook.sellHeap.AddOrder(12, 50)
	orderBook.sellHeap.AddOrder(12, 70)
	orderBook.sellHeap.AddOrder(11, 75)
	orderBook.sellHeap.AddOrder(5, 75)

	// cancel a buy order
	if err := orderBook.buyHeap.CancelOrder(buyOrderID); err != nil {
		fmt.Printf("Error cancelling buy order: %v\n", err)
	}

	// cancel a sell order
	if err := orderBook.sellHeap.CancelOrder(sellOrderID1); err != nil {
		fmt.Printf("Error cancelling sell order: %v\n", err)
	}

	// match orders in the order book
	orderBook.MatchOrders()

	// print out the remaining orders in the order book
	fmt.Println("Buy orders:")
	for _, order := range orderBook.buyHeap.orders {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
	fmt.Println("Sell orders:")
	for _, order := range orderBook.sellHeap.orders {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
}
