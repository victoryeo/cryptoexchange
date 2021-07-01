package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
)

func main() {
	log.Printf("Start order matching\n")
	const ORDERBOOK_LEN = 100
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Printf("%v\n", rdb)

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, ORDERBOOK_LEN),
		SellOrders: make([]engine.Order, 0, ORDERBOOK_LEN),
	}
	log.Printf("%v\n", book)

	// create a signal channel to know when we are done
	done := make(chan bool)

	// start processing orders
	go func() {
		for i := 1; i < ORDERBOOK_LEN; i++ {
			var order engine.Order
			key := "order" + strconv.Itoa(i)
			allVal, err := rdb.HGetAll(ctx, key).Result()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", allVal)
			if len(allVal) == 0 {
				done <- true
			}

			//scanVal := rdb.HScan(ctx, "order", 0, "*", 0)
			//fmt.Println("scan", scanVal)

			for index, element := range allVal {
				fmt.Println("Key:", index, "=>", "Value:", element)
			}
			/*
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				// encoding the map
				err = enc.Encode(allVal)
				fmt.Printf("%v\n", buf)
				// convert to byte array
				var orderHolder []byte
				orderHolder = buf.Bytes()
				fmt.Printf("orderHolder\n")
				fmt.Printf("%v\n", orderHolder) */

			// decode the message
			price, err := strconv.Atoi(allVal["price"])
			order.Price = uint64(price)
			qty, err := strconv.Atoi(allVal["qty"])
			order.Quantity = uint32(qty)
			if allVal["type"] == "buy" {
				order.Type = 1
			} else {
				order.Type = 0
			}
			fmt.Printf("order\n")
			fmt.Printf("%v\n", order)
			// process the order
			//trades := book.Process(order)
			//fmt.Printf("%v\n", trades)
		}
	}()

	// wait until we are done
	<-done
}
