package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
)

func main() {
	log.Printf("Start order matching\n")
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Printf("%v\n", rdb)

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}
	log.Printf("%v\n", book)

	// create a signal channel to know when we are done
	done := make(chan bool)

	// start processing orders
	go func() {
		for {
			var order engine.Order
			val, err := rdb.HGetAll(ctx, "order").Result()
			if err != nil {
				//panic(err)
				done <- true
			}
			fmt.Printf("%v\n", val)

			scanVal := rdb.HScan(ctx, "order", 0, "*", 0)
			fmt.Println("scan", scanVal)

			for key, element := range val {
				fmt.Println("Key:", key, "=>", "Value:", element)
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			// Encoding the map
			err = enc.Encode(val)
			fmt.Printf("%v\n", buf)
			// convert to byte array
			var orderHolder []byte
			orderHolder = buf.Bytes()
			fmt.Printf("orderHolder\n")
			fmt.Printf("%v\n", orderHolder)

			// decode the message
			order.FromJSON(orderHolder)
			fmt.Printf("order\n")
			fmt.Printf("%d\n", order.Price)
			// process the order
			//trades := book.Process(order)
			//fmt.Printf("%v\n", trades)
		}
	}()

	// wait until we are done
	<-done
}
