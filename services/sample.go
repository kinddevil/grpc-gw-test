package services

import (
	"context"
	"google.golang.org/grpc/metadata"
	pb "grpc_tpl/service_interfaces"
	"log"
)

type Server struct{}

func (s *Server) Sample(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md)
	return &pb.Reply{Status: 1, Code: "ok", Msg: ""}, nil
}
