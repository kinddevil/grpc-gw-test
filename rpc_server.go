package main

import (
	"google.golang.org/grpc"
	pb "grpc_tpl/service_interfaces"
	"grpc_tpl/services"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	log.Printf("start Score gRPC server with %v", port)
	pb.RegisterSampleServiceServer(s, &services.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
