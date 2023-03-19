package main

import (
	"container/heap"
	"fmt"

	"github.com/google/uuid"
)

type Order struct {
	orderID  string
	price    int
	quantity int
}

type BuyHeap []Order

func (bh BuyHeap) Len() int           { return len(bh) }
func (bh BuyHeap) Less(i, j int) bool { return bh[i].price > bh[j].price } // use greater than to implement max heap
func (bh BuyHeap) Swap(i, j int)      { bh[i], bh[j] = bh[j], bh[i] }

func (bh *BuyHeap) Push(x interface{}) {
	*bh = append(*bh, x.(Order))
}

func (bh *BuyHeap) Pop() interface{} {
	old := *bh
	n := len(old)
	x := old[n-1]
	*bh = old[:n-1]
	return x
}

type SellHeap []Order

func (sh SellHeap) Len() int           { return len(sh) }
func (sh SellHeap) Less(i, j int) bool { return sh[i].price < sh[j].price } // use less than to implement min heap
func (sh SellHeap) Swap(i, j int)      { sh[i], sh[j] = sh[j], sh[i] }

func (sh *SellHeap) Push(x interface{}) {
	*sh = append(*sh, x.(Order))
}

func (sh *SellHeap) Pop() interface{} {
	old := *sh
	n := len(old)
	x := old[n-1]
	*sh = old[:n-1]
	return x
}

type OrderBook struct {
	buyHeap  BuyHeap
	sellHeap SellHeap
}

func (ob *OrderBook) AddBuyOrder(price, quantity int) string {
	orderID := uuid.New().String()
	heap.Push(&ob.buyHeap, Order{orderID, price, quantity})
	return orderID
}

func (ob *OrderBook) AddSellOrder(price, quantity int) string {
	orderID := uuid.New().String()
	heap.Push(&ob.sellHeap, Order{orderID, price, quantity})
	return orderID
}

func (ob *OrderBook) CancelBuyOrder(orderID string) {
	buyOrders := make([]Order, 0)
	for len(ob.buyHeap) > 0 {
		order := heap.Pop(&ob.buyHeap).(Order)
		if order.orderID == orderID {
			// found the order to cancel
			break
		} else {
			buyOrders = append(buyOrders, order)
		}
	}
	// add back the remaining buy orders to the buy heap
	for _, order := range buyOrders {
		heap.Push(&ob.buyHeap, order)
	}
}

func (ob *OrderBook) CancelSellOrder(orderID string) {
	sellOrders := make([]Order, 0)
	for len(ob.sellHeap) > 0 {
		order := heap.Pop(&ob.sellHeap).(Order)
		if order.orderID == orderID {
			// found the order to cancel
			break
		} else {
			sellOrders = append(sellOrders, order)
		}
	}
	// add back the remaining sell orders to the sell heap
	for _, order := range sellOrders {
		heap.Push(&ob.sellHeap, order)
	}
}

func (ob *OrderBook) MatchOrders() {
	for len(ob.buyHeap) > 0 && len(ob.sellHeap) > 0 {
		bestBuy := ob.buyHeap[0]
		bestSell := ob.sellHeap[0]
		if bestBuy.price >= bestSell.price {
			tradePrice := bestSell.price
			tradeQuantity := min(bestBuy.quantity, bestSell.quantity) // match the lowest quantity between the two orders
			fmt.Printf("Trade executed at price %d for quantity %d\n", tradePrice, tradeQuantity)
			if bestBuy.quantity > tradeQuantity {
				ob.buyHeap[0] = Order{bestBuy.orderID, bestBuy.price, bestBuy.quantity - tradeQuantity}
				heap.Fix(&ob.buyHeap, 0)
			} else {
				heap.Pop(&ob.buyHeap)
			}
			if bestSell.quantity > tradeQuantity {
				ob.sellHeap[0] = Order{bestSell.orderID, bestSell.price, bestSell.quantity - tradeQuantity}
				heap.Fix(&ob.sellHeap, 0)
			} else {
				heap.Pop(&ob.sellHeap)
			}
		} else {
			break
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// create a new order book
	orderBook := OrderBook{}

	// add buy and sell orders to the order book
	buyOrderID := orderBook.AddBuyOrder(10, 100)
	orderBook.AddBuyOrder(12, 100)
	orderBook.AddBuyOrder(11, 120)
	sellOrderID1 := orderBook.AddSellOrder(12, 50)
	orderBook.AddSellOrder(12, 70)
	orderBook.AddSellOrder(11, 75)

	// cancel a buy order
	orderBook.CancelBuyOrder(buyOrderID)

	// cancel a sell order
	orderBook.CancelSellOrder(sellOrderID1)

	// match orders in the order book
	orderBook.MatchOrders()

	// print out the remaining orders in the order book
	fmt.Println("Buy orders:")
	for _, order := range orderBook.buyHeap {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
	fmt.Println("Sell orders:")
	for _, order := range orderBook.sellHeap {
		fmt.Printf("%s: %d @ %d\n", order.orderID, order.quantity, order.price)
	}
}
