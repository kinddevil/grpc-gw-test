package services

import (
	"context"
	pb "github.com/kinddevil/grpc-gw-test/service_interfaces"
	"google.golang.org/grpc/metadata"
	"log"
)

type Server struct{}

func (s *Server) Sample(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md)
	return &pb.Reply{Status: 1, Code: "ok", Msg: ""}, nil
}
