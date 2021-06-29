package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"flag"
	"github.com/victoryeo/cryptoexchange/engine"
	pb "simple"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 19000, "The server port")	
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
	// write order to db
	// stub code

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
	log.Printf("Start order server\n")

	// create the order book
	log.Printf("Create order book\n")
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}
	fmt.Printf("%v\n",book)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Start grpc server\n")
	grpcServer := grpc.NewServer()

	log.Printf("Register trade engine server\n")
	pb.RegisterTradeEngineServer(grpcServer, &server{})

	log.Printf("Serve grpc\n")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}