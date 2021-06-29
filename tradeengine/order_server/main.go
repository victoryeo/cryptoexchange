package main

import (
	"log"
	"net"
	"github.com/victoryeo/cryptoexchange/engine"
	pb "simple"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// implement Server.
type server struct {
	pb.UnimplementedTradeEngineServer
}

// implements rpc
func (s *server) GetOrder(in *pb.Empty, stream pb.TradeEngine_GetOrderServer) error {
	return nil
}
func (s *server) GetTrade(in *pb.Empty, stream pb.TradeEngine_GetTradeServer) error {
	return nil
}
func (s *server) SendTrade(stream pb.TradeEngine_SendTradeServer) error {
	return nil
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTradeEngineServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}

	// create a signal channel to know when we are done
	done := make(chan bool)

	// start processing orders
	go func() {
		for msg := range consumer.Messages() {
			var order engine.Order
			// decode the message
			order.FromJSON(msg.Value)
			// process the order
			trades := book.Process(order)
			// send trades to message queue
			for _, trade := range trades {
				rawTrade := trade.ToJSON()
				
			}
			// mark the message as processed
			consumer.MarkOffset(msg, "")
		}
		done <- true
	}()

	// wait until we are done
	<-done
}