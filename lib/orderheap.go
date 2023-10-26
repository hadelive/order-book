package orderBook

const (
	Buy  = "Buy"
	Sell = "Sell"
)

type Order struct {
	OrderID  string
	Side     string
	Price    int64
	Quantity int64
}

// OrderHeap implements a heap for orders, with the ordering determined by the Compare function.
type OrderHeap struct {
	orders []*Order
	less   func(o1, o2 *Order) bool
}

func (h OrderHeap) Len() int           { return len(h.orders) }
func (h OrderHeap) Less(i, j int) bool { return h.less(h.orders[i], h.orders[j]) }
func (h OrderHeap) Swap(i, j int)      { h.orders[i], h.orders[j] = h.orders[j], h.orders[i] }

func (h *OrderHeap) Push(x interface{}) {
	h.orders = append(h.orders, x.(*Order))
}

func (h *OrderHeap) Pop() interface{} {
	old := h.orders
	n := len(old)
	x := old[n-1]
	h.orders = old[0 : n-1]
	return x
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
