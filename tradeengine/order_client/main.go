package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "simple"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:19000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	log.Printf("Start order client with tls %t\n", *tls)
	if *tls {
		if *caFile == "" {
			//to be added
			*caFile = ""
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	log.Printf("Connect to server\n")
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTradeEngineClient(conn)

	log.Printf("Call SendOrder\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	order := &pb.Order{Id: 1, Price: 10, Quantity: 100, Type: "buy"}
	ret, err := client.SendOrder(ctx, order)
	log.Printf("SendOrder return value %d\n", ret)

	order = &pb.Order{Id: 2, Price: 10, Quantity: 100, Type: "sell"}
	ret, err = client.SendOrder(ctx, order)
	log.Printf("SendOrder return value %d\n", ret)
}
