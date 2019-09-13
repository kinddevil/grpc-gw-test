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

var (
	grpcPort    = cfgs.GetString("grpc.port")
	maxConnIdle = cfgs.GetInt("grpc.max_connection_idle")
	timeOut     = cfgs.GetInt("grpc.time_out")
)

func ServeGRPC(terminate chan<- CancelFun) {
	lis, err := net.Listen("tcp", grpcPort)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(maxConnIdle) * time.Second,
			Timeout:           time.Duration(timeOut) * time.Second,
		}),
		grpc.UnaryInterceptor(
			inspectors.MiddlewareFunc(inspectors.GetUserInfo),
		),
	)

	terminate <- func() error {
		s.Stop()
		return nil
	}

	log.Printf("start gRPC server with %v", grpcPort)
	pb.RegisterSampleServiceServer(s, &services.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
