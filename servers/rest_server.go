package servers

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"grpc_tpl/configs"
	pb "grpc_tpl/service_interfaces"
	"log"
	"net/http"
)

var (
	cfgs               = configs.CONFIGS
	//grpcPort           = cfgs.GetString("grpc.port")
	//grpcServerEndpoint = cfgs.GetString("rest.grpc_addr")
)

func ServeHttp(terminate chan<- func() error) {
	log.Println(cfgs)
	grpcPort           := cfgs.GetString("grpc.port")
	grpcServerEndpoint := cfgs.GetString("rest.grpc_addr")


	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	port := grpcPort

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	s := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// TODO pass the parameters like endpoint and opts
	err := pb.RegisterSampleServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	terminate <- func() error {
		return s.Shutdown(ctx)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("start service with %v", port)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
