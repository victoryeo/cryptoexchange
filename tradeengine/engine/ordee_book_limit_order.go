package engine

import "fmt"

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order Order) []Trade {
	if order.Type == 1 {
		return book.processLimitBuy(order)
	} else if order.Type == 0 {
		return book.processLimitSell(order)
	} else {
		fmt.Print("wrong order type\n")
		return nil
	}
}

// Process a limit buy order
func (book *OrderBook) processLimitBuy(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(book.SellOrders)
	fmt.Printf("slen %d\n", n)
	// check if we have at least one matching order
	if n >= 1 {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			// break from for loop if sell price in Q is > order price
			if sellOrder.Price > order.Price {
				break
			}
			// fill the entire order
			if sellOrder.Quantity >= order.Quantity {
				trades = append(trades, Trade{order.Id, sellOrder.Id, order.Quantity, sellOrder.Price})
				sellOrder.Quantity -= order.Quantity
				if sellOrder.Quantity == 0 {
					book.removeSellOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if sellOrder.Quantity < order.Quantity {
				trades = append(trades, Trade{order.Id, sellOrder.Id, sellOrder.Quantity, sellOrder.Price})
				order.Quantity -= sellOrder.Quantity
				book.removeSellOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addBuyOrder(order)
	return trades
}

// Process a limit sell order
func (book *OrderBook) processLimitSell(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	n := len(book.BuyOrders)
	fmt.Printf("blen %d\n", n)
	// check if we have at least one matching order
	if n >= 1 {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			// break from for loop if buy price in Q is < order price
			if buyOrder.Price < order.Price {
				break
			}
			// fill the entire order
			if buyOrder.Quantity >= order.Quantity {
				trades = append(trades, Trade{order.Id, buyOrder.Id, order.Quantity, buyOrder.Price})
				buyOrder.Quantity -= order.Quantity
				if buyOrder.Quantity == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if buyOrder.Quantity < order.Quantity {
				trades = append(trades, Trade{order.Id, buyOrder.Id, buyOrder.Quantity, buyOrder.Price})
				order.Quantity -= buyOrder.Quantity
				book.removeBuyOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addSellOrder(order)
	return trades
}
