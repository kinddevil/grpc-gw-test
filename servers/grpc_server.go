package servers

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"grpc_tpl/inspectors"
	pb "grpc_tpl/service_interfaces"
	"grpc_tpl/services"
	"log"
	"net"
	"time"
)

const (
	port = ":50051"
)

func ServeGRPC(terminate chan<- CancelFun) {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
		grpc.UnaryInterceptor(
			inspectors.MiddlewareFunc(inspectors.GetUserInfo),
		),
	)

	terminate <- func() error {
		s.Stop()
		return nil
	}

	log.Printf("start gRPC server with %v", port)
	pb.RegisterSampleServiceServer(s, &services.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
