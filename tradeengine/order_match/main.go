package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
)

type Response struct {
	Title string `json:"title"`
}

func crypto_interface() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("http do error ")
		fmt.Println(err.Error())
	}

	defer func() {
		if resp != nil {
			err = resp.Body.Close()
			if err != nil {
				fmt.Print("http close error ")
				fmt.Print(err.Error())
			}
		}
	}()
	if resp != nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print("http read error ")
			fmt.Print(err.Error())
		}
		var responseObject Response
		json.Unmarshal(bodyBytes, &responseObject)
		fmt.Printf("API Response as struct %+v\n", responseObject)
	}
}

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
			fmt.Printf("allval %d\n", len(allVal))
			if len(allVal) == 0 {
				done <- true
			} else {
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
				if err != nil {
					log.Println(err)
				}
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
				trades := book.Process(order)
				fmt.Printf("match trade %v\n", trades)

				// forware the trade to crypto API
				crypto_interface()
			}
		}
	}()

	// wait until we are done
	<-done

	fmt.Printf("%v\n", book)
}
