package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	pb "grpc_tpl/service_interfaces"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func startHttpServer() {
	ctx := context.Background()
	// TODO config the timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		//grpc.WithDisableRetry(),
		//grpc.WithTimeout(10 * time.Second),
	}

	// TODO pass the parameters like endpoint and opts
	err := pb.RegisterSampleServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	// TODO config the port
	port := ":8081"

	s := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("start service with %v", port)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func run() {
	cancelFunChan := make(chan context.CancelFunc, 1)
	log.Println("before")
	log.Println(cancelFunChan)
	//go startHttpServer()
	startHttpServer()
	cancelFun := <-cancelFunChan
	time.Sleep(3 * time.Second)
	log.Println("after")
	log.Println(cancelFunChan)
	log.Println(cancelFun)
	cancelFun()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	log.Printf("Stop HTTP service with sig %v", sig)
}

func main() {
	run()
	log.Println("The server is stopped")
}
