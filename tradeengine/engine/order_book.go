package engine

import "fmt"

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

type MapOrderBook map[string]OrderBook

// Add a buy order to the order book
func (book *OrderBook) addBuyOrder(order Order) {
	n := len(book.BuyOrders)
	fmt.Printf("%d\n", n)
	if n == 0 {
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		var i int
		for i = n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price <= order.Price {
				break
			}
		}
		if i == n-1 {
			book.BuyOrders = append([]Order{order}, book.BuyOrders...)
		} else if i == n-1 {
			if book.BuyOrders[i].Price == order.Price {
				fmt.Printf("%f %d\n", book.BuyOrders[i].Price, i)
				book.BuyOrders = append(book.BuyOrders[:i+1], book.BuyOrders[i:]...)
				book.BuyOrders[i] = order
			} else {
				book.BuyOrders = append(book.BuyOrders, order)
			}
		} else {
			i++
			book.BuyOrders = append(book.BuyOrders[:i+1], book.BuyOrders[i:]...)
			book.BuyOrders[i] = order
		}
	}
}

// Add a sell order to the order book
func (book *OrderBook) addSellOrder(order Order) {
	n := len(book.SellOrders)
	if n == 0 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		var i int
		for i = n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price >= order.Price {
				break
			}
		}
		if i == n-1 {
			book.SellOrders = append([]Order{order}, book.SellOrders...)
		} else if i == n-1 {
			if book.SellOrders[i].Price == order.Price {
				fmt.Printf("%f %d\n", book.SellOrders[i].Price, i)
				book.SellOrders = append(book.SellOrders[:i+1], book.SellOrders[i:]...)
				book.SellOrders[i] = order
			} else {
				fmt.Printf("%f\n", book.SellOrders[i].Price)
				book.SellOrders = append(book.SellOrders, order)
			}
		} else {
			i++
			book.SellOrders = append(book.SellOrders[:i+1], book.SellOrders[i:]...)
			book.SellOrders[i] = order
		}
	}
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
