package orderBook

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Order represents an order in the order book.
type Order struct {
	orderID  string
	price    int
	quantity int
}

// OrderHeap implements a heap for orders, with the ordering determined by the Compare function.
type OrderHeap struct {
	orders []Order
	less   func(o1, o2 Order) bool
}

func (h OrderHeap) Len() int           { return len(h.orders) }
func (h OrderHeap) Less(i, j int) bool { return h.less(h.orders[i], h.orders[j]) }
func (h OrderHeap) Swap(i, j int)      { h.orders[i], h.orders[j] = h.orders[j], h.orders[i] }

func (h *OrderHeap) Push(x interface{}) {
	h.orders = append(h.orders, x.(Order))
}

func (h *OrderHeap) Pop() interface{} {
	old := h.orders
	n := len(old)
	x := old[n-1]
	h.orders = old[0 : n-1]
	return x
}

// OrderBook represents the order book, with a buy heap and a sell heap.
type OrderBook struct {
	BuyHeap  OrderHeap
	SellHeap OrderHeap
}

func NewOrderBook() OrderBook {
	return OrderBook{
		BuyHeap: OrderHeap{
			less: func(o1, o2 Order) bool { return o1.price > o2.price }, // use greater than for buy orders
		},
		SellHeap: OrderHeap{
			less: func(o1, o2 Order) bool { return o1.price < o2.price }, // use less than for sell orders
		},
	}
}

// AddOrder adds an order to the heap and returns the order ID.
func (h *OrderHeap) AddOrder(price, quantity int) string {
	orderID := uuid.New().String()
	heap.Push(h, Order{orderID, price, quantity})
	return orderID
}

// CancelOrder cancels an order with the given order ID from the heap.
// Returns an error if the order is not found.
func (h *OrderHeap) CancelOrder(orderID string) error {
	found := false
	orders := make([]Order, 0)
	for len(h.orders) > 0 {
		order := heap.Pop(h).(Order)
		if order.orderID == orderID {
			// found the order to cancel
			found = true
			fmt.Println("Cancelling order with ID: ", order.orderID)
			break
		} else {
			orders = append(orders, order)
		}
	}
	if !found {
		return errors.New("order not found")
	}
	// add back the remaining orders to the heap
	for _, order := range orders {
		heap.Push(h, order)
	}
	return nil
}

// MatchOrders matches the buy and sell orders in the order book and executes trades.
func (ob *OrderBook) MatchOrders() {
	for len(ob.BuyHeap.orders) > 0 && len(ob.SellHeap.orders) > 0 {
		bestBuy := ob.BuyHeap.orders[0]
		bestSell := ob.SellHeap.orders[0]
		if bestBuy.price >= bestSell.price {
			tradePrice := bestSell.price
			tradeQuantity := min(bestBuy.quantity, bestSell.quantity) // match the lowest quantity between the two orders
			fmt.Printf("Best buy quantity %d\n", bestBuy.quantity)
			fmt.Printf("Best sell quantity %d\n", bestSell.quantity)
			fmt.Printf("Trade executed at price %d for quantity %d\n", tradePrice, tradeQuantity)
			fmt.Printf("-----------------------------------------------------------\n")
			if bestBuy.quantity > tradeQuantity {
				ob.BuyHeap.orders[0] = Order{bestBuy.orderID, bestBuy.price, bestBuy.quantity - tradeQuantity}
				heap.Fix(&ob.BuyHeap, 0)
			} else {
				heap.Pop(&ob.BuyHeap)
			}
			if bestSell.quantity > tradeQuantity {
				ob.SellHeap.orders[0] = Order{bestSell.orderID, bestSell.price, bestSell.quantity - tradeQuantity}
				heap.Fix(&ob.SellHeap, 0)
			} else {
				heap.Pop(&ob.SellHeap)
			}
		} else {
			break
		}
	}
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
