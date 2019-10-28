package servers

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"grpc-gw-test/cluster"
	"grpc-gw-test/inspectors"
	pb "grpc-gw-test/service_interfaces"
	"grpc-gw-test/services"
	"log"
	"net"
	"time"
)

const (
	LB_TTL = 10 // seconds
)

func ServeGRPC(terminate chan<- CancelFun, cfgs *viper.Viper) {
	grpcHost := cfgs.GetString("grpc.host") // xx.xx.xx.xx
	grpcPort := cfgs.GetString("grpc.port") // :xxx
	addr := fmt.Sprintf("%v%v", grpcHost, grpcPort)

	maxConnIdle := cfgs.GetInt("grpc.max_connection_idle")
	timeOut := cfgs.GetInt("grpc.time_out")

	etcdAddr := cfgs.GetString("common.etcd_addrs")
	servName := cfgs.GetString("common.service_name")

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

	// Register service
	registerServer := &cluster.RegisterService{}
	registerServer.Register(etcdAddr, servName, addr, LB_TTL)

	terminate <- func() error {
		registerServer.UnRegister(servName, addr)
		s.Stop()
		return nil
	}

	log.Printf("start gRPC server with %v", grpcPort)
	pb.RegisterSampleServiceServer(s, &services.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
