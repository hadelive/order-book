package orderBook

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
)

func TestOrderBook(t *testing.T) {
	ob := NewOrderBook()

	// Add some orders
	order1 := &Order{OrderID: "1", Side: Buy, Price: 100, Quantity: 10}
	order2 := &Order{OrderID: "2", Side: Buy, Price: 90, Quantity: 5}
	order3 := &Order{OrderID: "3", Side: Sell, Price: 110, Quantity: 7}
	order4 := &Order{OrderID: "4", Side: Sell, Price: 120, Quantity: 3}
	order5 := &Order{OrderID: "5", Side: Buy, Price: 120, Quantity: 11}

	err := ob.AddOrder(order1)
	if err != nil {
		t.Errorf("Error adding order1: %s", err)
	}
	// Validate the state of the order book after canceling the order
	expectedBuyOrders := []*Order{
		{OrderID: "1", Side: Buy, Price: 100, Quantity: 10},
	}

	validateOrderBookState(t, ob.BuyHeap.orders, expectedBuyOrders)
	validateOrderBookState(t, ob.SellHeap.orders, []*Order{})
	err = ob.AddOrder(order2)
	if err != nil {
		t.Errorf("Error adding order2: %s", err)
	}
	expectedBuyOrders = []*Order{
		{OrderID: "1", Side: Buy, Price: 100, Quantity: 10},
		{OrderID: "2", Side: Buy, Price: 90, Quantity: 5},
	}
	validateOrderBookState(t, ob.BuyHeap.orders, expectedBuyOrders)
	validateOrderBookState(t, ob.SellHeap.orders, []*Order{})
	// Cancel an order
	err = ob.CancelOrder(order2)
	if err != nil {
		t.Errorf("Error canceling order: %s", err)
	}

	expectedBuyOrders = []*Order{
		{OrderID: "1", Side: Buy, Price: 100, Quantity: 10},
	}
	validateOrderBookState(t, ob.BuyHeap.orders, expectedBuyOrders)
	validateOrderBookState(t, ob.SellHeap.orders, []*Order{})

	fmt.Println("Add sell orders:")
	err = ob.AddOrder(order3)
	if err != nil {
		t.Errorf("Error adding order3: %s", err)
	}

	err = ob.AddOrder(order4)
	if err != nil {
		t.Errorf("Error adding order4: %s", err)
	}

	expectedSellOrders := []*Order{
		{OrderID: "3", Side: Sell, Price: 110, Quantity: 7},
		{OrderID: "4", Side: Sell, Price: 120, Quantity: 3},
	}

	validateOrderBookState(t, ob.SellHeap.orders, expectedSellOrders)

	err = ob.AddOrder(order5)
	if err != nil {
		t.Errorf("Error adding order5: %s", err)
	}

	expectedBuyOrders = []*Order{
		{OrderID: "5", Side: Buy, Price: 120, Quantity: 1},
		{OrderID: "1", Side: Buy, Price: 100, Quantity: 10},
	}
	validateOrderBookState(t, ob.BuyHeap.orders, expectedBuyOrders)
	validateOrderBookState(t, ob.SellHeap.orders, []*Order{})
}

func validateOrderBookState(t *testing.T, orders []*Order, expectedOrders []*Order) {
	fmt.Println("validate:")
	for _, order := range orders {
		fmt.Println("orders: ", order.OrderID, ", ", order.Side, ", ", order.Price, ", ", order.Quantity)
	}
	if len(orders) != len(expectedOrders) {
		t.Errorf("Expected %d orders, got %d", len(expectedOrders), len(orders))
		return
	}

	for i, order := range orders {
		expectedOrder := expectedOrders[i]
		if order.OrderID != expectedOrder.OrderID && order.Side != expectedOrder.Side && order.Quantity != expectedOrder.Quantity && order.Price != expectedOrder.Price {
			t.Errorf("Order mismatch at index %d. Expected %v, got %v", i, expectedOrder, order)
		}
	}
}

func BenchmarkAddOrder(b *testing.B) {
	orderCount := []int{1, 100, 1000, 10000, 100000, 1000000}

	for _, count := range orderCount {
		b.Run(fmt.Sprintf("OrderCount_%d", count), func(b *testing.B) {
			ob := NewOrderBook()
			orders := generateOrders(count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, order := range orders {
					_ = ob.AddOrder(order)
				}
			}
		})
	}
}

func BenchmarkCancelOrder(b *testing.B) {
	orderCount := []int{1, 100, 1000, 10000, 100000}

	for _, count := range orderCount {
		b.Run(fmt.Sprintf("OrderCount_%d", count), func(b *testing.B) {
			ob := NewOrderBook()
			orders := generateOrders(count)

			for _, order := range orders {
				_ = ob.AddOrder(order)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, order := range orders {
					_ = ob.CancelOrder(order)
				}
			}
		})
	}
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	basePrice := decimal.NewFromFloat(1e4)
	for i := 0; i < count; i++ {
		priceFluctuation := decimal.NewFromFloat(rand.Float64()*200 - 100)
		price := basePrice.Add(priceFluctuation).Round(2)

		sizeFluctuation := decimal.NewFromFloat(rand.Float64()*10 - 5)
		size := decimal.NewFromFloat(10).Add(sizeFluctuation).Round(2)

		orders[i] = &Order{
			OrderID:  strconv.Itoa(i + 1),
			Side:     Buy,
			Price:    price.IntPart(),
			Quantity: size.IntPart(),
		}
	}
	return orders
}
