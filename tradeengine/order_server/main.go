package main

import (
	"context"
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
	sob []*pb.OrderBook
	gob []engine.OrderBook
}

// implements rpc
func (s *server) SendOrder(ctx context.Context, msg *pb.Order) (*pb.Empty, error) {
	var data engine.Order
	data.Price = msg.Price
	data.Amount = msg.Quantity
	if msg.Type == "buy" {
		data.Side = 1
	} else {
		data.Side = 0
	}
	//write order to db
	return &pb.Empty{}, nil
}

func (s *server) GetOrderStream(in *pb.Empty, stream pb.TradeEngine_GetOrderStreamServer) error {

	return nil
}
func (s *server) GetTradeStream(in *pb.Empty, stream pb.TradeEngine_GetTradeStreamServer) error {
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
		for {
			// get order from db

			// process the order
			trades := book.Process(order)
			// mark the message as processed
		}
		done <- true
	}()

	// wait until we are done
	<-done
}