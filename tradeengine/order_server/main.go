package main

import (
	"example.com/engine"
	"log"
	pb "simple"

)

const (
	port = ":50051"
)

// implement Server.
type server struct {
	pb.UnimplementedTradeEngineServer
}

// implements rpc
func (s *server) GetOrder(ctx context.Context, in *pb.OrderID) (*pb.OrderMessage, error) {
	log.Printf("Received: %v", in.GetName())
	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
func (s *server) GetTrade(ctx context.Context, in *pb.OrderID) (*pb.TradeMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrade not implemented")
}
func (s *server) SendTrade(ctx context.Context, in *pb.TradeMessage) (*pb.OrderID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTrade not implemented")
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
				producer.Input() <- &sarama.ProducerMessage{
					Topic: "trades",
					Value: sarama.ByteEncoder(rawTrade),
				}
			}
			// mark the message as processed
			consumer.MarkOffset(msg, "")
		}
		done <- true
	}()

	// wait until we are done
	<-done
}