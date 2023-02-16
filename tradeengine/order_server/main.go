package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "simple"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/victoryeo/cryptoexchange/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 19000, "The server port")
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
	data.Quantity = msg.Quantity
	data.Id = strconv.FormatInt(int64(msg.Id), 10)

	log.Printf("Order id %d received\n", msg.Id)
	if msg.Type == "buy" {
		data.Type = 1
	} else {
		data.Type = 0
	}
	// write order to db
	// redis code
	//key := "order" + strconv.Itoa(int(msg.Id))
	key := msg.TokenName + msg.TokenType
	err := s.rdb.HSet(ctx, key, "id", msg.Id,
		"price", msg.Price, "qty", msg.Quantity,
		"type", msg.Type, "tokenName", msg.TokenName,
		"tokenType", msg.TokenType, "processed", false).Err()
	if err != nil {
		panic(err)
	}
	val, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("order", val)

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

func initServer() *teServer {
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
	log.Printf("Start order server with tls %t\n", *tls)

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
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = ""
		}
		if *keyFile == "" {
			*keyFile = ""
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	log.Printf("Start grpc server\n")
	grpcServer := grpc.NewServer(opts...)

	log.Printf("Register trade engine server\n")
	pb.RegisterTradeEngineServer(grpcServer, initServer())

	log.Printf("Serve grpc\n")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
