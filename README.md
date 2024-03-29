# Order Books

This is a simple implementation of an order book using two heaps - a max heap for buy orders and a min heap for sell orders. The order book allows users to add buy or sell orders, cancel existing orders, and match buy and sell orders to execute trades.

## Why Use Heaps?

Heaps are a natural choice for implementing an order book, as they allow us to efficiently retrieve the best buy and sell orders from the book. The max heap for buy orders ensures that the highest buy order (i.e. the highest bidding Price) is always on top, while the min heap for sell orders ensures that the lowest sell order (i.e. the lowest asking Price) is always on top. In addition, heaps have a logarithmic time complexity for both insertion and removal of elements, which makes them a fast and efficient data structure for handling large numbers of orders in real time.

## How to Use

To use the order book, simply create a new instance of the OrderBook struct and call its methods to add, cancel or match orders. The AddBuyOrder and AddSellOrder methods add new orders to the order book and return a unique order ID, while the CancelBuyOrder and CancelSellOrder methods cancel an existing order with the given ID. The MatchOrders method matches the best buy and sell orders in the order book and executes trades.

### Test

```sh
go test
```

## TODO

- [ ] Save/restore the state of the order book
- [ ] Handle buy and sell events
- [ ] Save the order history

## References

- <https://www.youtube.com/watch?v=dUMWMZmMsVE&ab_channel=TechDummiesNarendraL>
