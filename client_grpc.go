package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_tpl/service_interfaces"
	"log"
	"time"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSampleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Sample(ctx, &pb.Request{Id: "1", Name: "anonymous"})
	if err != nil {
		log.Fatalf("could not score: %v", err)
	}
	log.Printf("Greeting: %v %v %v", r.Status, r.Code, r.Msg)
}
