package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
)

type Response struct {
	Title string `json:"title"`
}

type SendResponse struct {
	Address string `json:"address"`
	Amount  uint32 `json:"amount"`
}

func crypto_interface_init() {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/", nil)
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

func crypto_interface_send(qty uint32) {
	client := &http.Client{}

	const address = "fe30945738"
	//var jsonStr = []byte(`{"amount":qty}`)
	pb := &SendResponse{Address: address, Amount: qty}
	jsonStr, err := json.Marshal(pb)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/send/"+address, bytes.NewBuffer(jsonStr))
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
		var responseObject SendResponse
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

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, ORDERBOOK_LEN),
		SellOrders: make([]engine.Order, 0, ORDERBOOK_LEN),
	}
	log.Printf("%v\n", book)

	// create a signal channel to know when we are done
	done := make(chan bool, 1)

	//var mOrderBook = make(engine.MapOrderBook, 0)

	// start processing orders
	go func() {
		for i := 1; i < ORDERBOOK_LEN; i++ {
			var order engine.Order
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				done <- false
			default:
				orderKey := "order" + strconv.Itoa(i)
				allVal, err := rdb.HGetAll(ctx, orderKey).Result()
				if err != nil {
					panic(err)
				}

				if len(allVal) == 0 {
					//done <- true
					time.Sleep(2 * time.Second)
					continue
				} else {
					fmt.Printf("%v\n", allVal)
					fmt.Printf("allval %d\n", len(allVal))
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
					if len(trades) > 0 {
						crypto_interface_init()
						crypto_interface_send(order.Quantity)
					}
				}
			}
		}
	}()

	// wait until we are done
	<-done

	fmt.Printf("%v\n", book)
}
