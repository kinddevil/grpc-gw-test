package servers

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	etcdAddr := cfgs.GetString("common.register_etcd_service")
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
			grpc_middleware.ChainUnaryServer(
				inspectors.MiddlewareFunc(inspectors.GetUserInfo),
				grpc_prometheus.UnaryServerInterceptor,
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_prometheus.StreamServerInterceptor,
			),
		),
	)

	// Register prometheus and open http port for scraping
	// TODO abstract metrics logic
	grpc_prometheus.Register(s)
	http.Handle("/metrics", promhttp.Handler())
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	customizedCounterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	// TODO config port
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()
	// TODO end

	// Register service
	registerServer := &cluster.RegisterService{}
	registerServer.Register(etcdAddr, servName, addr, LB_TTL)

	terminate <- func() error {
		registerServer.UnRegister(servName, addr)
		s.Stop()
		httpServer.Close()
		return nil
	}

	log.Printf("start gRPC server with %v", grpcPort)
	pb.RegisterSampleServiceServer(s, &services.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
