package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "simple"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 19000, "The server port")
)

// implement Server.
type teServer struct {
	pb.UnimplementedTradeEngineServer
	sob []*pb.OrderBook
	gob []engine.OrderBook
	rdb *redis.Client
}

// implements rpc
func (s *teServer) SendOrder(ctx context.Context, msg *pb.Order) (*pb.Empty, error) {
	var data engine.Order
	data.Price = msg.Price
	data.Amount = msg.Quantity

	log.Printf("Order id %d received\n", msg.Id)
	if msg.Type == "buy" {
		data.Side = 1
	} else {
		data.Side = 0
	}
	// write order to db
	// stub code
	err := s.rdb.HSet(ctx, "order", "id", msg.Id, "price", msg.Price).Err()
	if err != nil {
		panic(err)
	}
	val, err := s.rdb.HGetAll(ctx, "order").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	return &pb.Empty{}, nil
}

func (s *teServer) GetOrderStream(in *pb.Empty, stream pb.TradeEngine_GetOrderStreamServer) error {

	return nil
}
func (s *teServer) GetTradeStream(in *pb.Empty, stream pb.TradeEngine_GetTradeStreamServer) error {
	return nil
}
func (s *teServer) SendTrade(stream pb.TradeEngine_SendTradeServer) error {
	return nil
}

func newServer() *teServer {
	//var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Printf("%v\n", rdb)
	s := &teServer{rdb: rdb}
	return s
}

func main() {
	log.Printf("Start order server\n")

	// create the order book
	log.Printf("Create order book\n")
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}
	fmt.Printf("%v\n", book)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Start grpc server\n")
	grpcServer := grpc.NewServer()

	log.Printf("Register trade engine server\n")
	pb.RegisterTradeEngineServer(grpcServer, newServer())

	log.Printf("Serve grpc\n")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
