package main

import (
	"context"
	pb "github.com/kinddevil/grpc-gw-test/service_interfaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

const (
	address = "localhost:50051"
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

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Meta-TrackId", "client-sent")

	r, err := c.Sample(ctx, &pb.Request{Id: "1", Name: "anonymous"})
	if err != nil {
		log.Fatalf("could not score: %v", err)
	}
	log.Printf("Greeting: %v %v %v", r.Status, r.Code, r.Msg)
}
