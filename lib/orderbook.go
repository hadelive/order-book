package orderBook

import (
	"container/heap"
	"errors"
	"fmt"
)

type OrderBook struct {
	BuyHeap  OrderHeap
	SellHeap OrderHeap
}

func NewOrderBook() OrderBook {
	return OrderBook{
		BuyHeap: OrderHeap{
			less: func(o1, o2 *Order) bool {
				if o1.Price == o2.Price {
					return o1.OrderID < o2.OrderID // prioritize older order if Prices are equal
				}
				return o1.Price > o2.Price // use greater than for buy orders
			},
		},
		SellHeap: OrderHeap{
			less: func(o1, o2 *Order) bool {
				if o1.Price == o2.Price {
					return o1.OrderID < o2.OrderID // prioritize older order if Prices are equal
				}
				return o1.Price < o2.Price // use less than for sell orders
			},
		},
	}
}

// AddOrder adds an order to the heap and returns the order ID.
func (ob *OrderBook) AddOrder(order *Order) error {
	if order.Quantity < 0 {
		return errors.New("quantity must be positive")
	}

	// Check if there are any matching orders in the opposite heap
	oppositeHeap := &ob.SellHeap
	if order.Side == Sell {
		oppositeHeap = &ob.BuyHeap
	}
	for len(oppositeHeap.orders) > 0 {
		oppositeOrder := oppositeHeap.orders[0]
		if (order.Side == Buy && order.Price < oppositeOrder.Price) || (order.Side == Sell && order.Price > oppositeOrder.Price) {
			break
		}
		// Execute the trade
		tradePrice := oppositeOrder.Price
		tradeQuantity := min(order.Quantity, oppositeOrder.Quantity)
		fmt.Printf("Trade executed at Price %d for Quantity %d\n", tradePrice, tradeQuantity)
		fmt.Printf("-----------------------------------------------------------\n")
		order.Quantity -= tradeQuantity
		oppositeOrder.Quantity -= tradeQuantity
		if oppositeOrder.Quantity == 0 {
			heap.Pop(oppositeHeap)
		}
		if order.Quantity == 0 {
			return nil
		}
	}
	if order.Quantity > 0 {
		if order.Side == Buy {
			heap.Push(&ob.BuyHeap, order)
		} else {
			heap.Push(&ob.SellHeap, order)
		}
	}
	return nil
}

// CancelOrder cancels the given order from the heap.
// Returns an error if the order is not found.
func (ob *OrderBook) CancelOrder(orderToCancel *Order) error {
	var heapToSearch *OrderHeap

	if orderToCancel.Side == Buy {
		heapToSearch = &ob.BuyHeap
	} else {
		heapToSearch = &ob.SellHeap
	}

	// Find the index of the order to cancel in the heap
	index := -1
	for i, order := range heapToSearch.orders {
		if order.OrderID == orderToCancel.OrderID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("order not found")
	}

	// Remove the order from the heap by swapping with the last order
	lastIndex := len(heapToSearch.orders) - 1
	heapToSearch.orders[index], heapToSearch.orders[lastIndex] = heapToSearch.orders[lastIndex], heapToSearch.orders[index]
	heapToSearch.orders = heapToSearch.orders[:lastIndex]

	// Heapify the heap after removing the order
	heap.Init(heapToSearch)

	return nil
}
