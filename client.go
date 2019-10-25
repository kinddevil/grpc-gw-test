package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
	pb "grpc-gw-test/service_interfaces"
	"time"
)
func main() {
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		log.Println("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSampleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Meta-TrackId", "client-sent")

	res, err := c.Sample(ctx, &pb.Request{Id: "1", Name: "anonymous"})
	log.Println(res)
	if err != nil {
		log.Println("could not score: %v", err)
	}
}
