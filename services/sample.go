package services

import (
	"context"
	pb "grpc_tpl/service_interfaces"
)

type Server struct{}

func (s *Server) Sample(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	return &pb.Reply{Status: 1, Code: "ok", Msg: ""}, nil
}
