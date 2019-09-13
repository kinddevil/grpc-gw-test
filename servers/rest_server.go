package servers

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	pb "grpc_tpl/service_interfaces"
	"log"
	"net/http"
)

var (
	restPort           = cfgs.GetString("rest.port")
	grpcServerEndpoint = cfgs.GetString("rest.grpc_addr")
)

func ServeHttp(terminate chan<- func() error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	s := &http.Server{
		Addr:    restPort,
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
	log.Printf("start service with %v", restPort)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
