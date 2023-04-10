package orderbook

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestOrderBook(t *testing.T) {
	// set the random seed to the current time
	rand.Seed(time.Now().UnixNano())
	// create a new order book
	orderBook := NewOrderBook()

	// sample buy/sell orders
	buyOrderID := orderBook.buyHeap.AddOrder(100, 100)
	sellOrderID := orderBook.sellHeap.AddOrder(110, 100)

	// create a wait group to wait for all Goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutines 1: continuosly add buy orders
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Println("Adding buy orders...", i)
			orderBook.buyHeap.AddOrder(100-i*10, i+1)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		}
	}()

	// Goroutines 2: continously add sell orders
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Println("Adding sell orders....", i)
			orderBook.sellHeap.AddOrder(10+i*10, i*2+1)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		}
	}()

	// wait for all Goroutines to finish
	wg.Wait()

	orderBook.buyHeap.CancelOrder(buyOrderID)
	orderBook.sellHeap.CancelOrder(sellOrderID)
	orderBook.MatchOrders()

	fmt.Println("Buy orders:")
	for _, order := range orderBook.buyHeap.orders {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
	fmt.Println("Sell orders:")
	for _, order := range orderBook.sellHeap.orders {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
}
