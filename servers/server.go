package servers

import (
	"grpc-gw-test/configs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type CancelFun = func() error
type Server func(chan<- func() error)

var (
	cfgs = configs.CONFIGS
)

func StartServers() {
	servers := []struct {
		name   string
		server Server
		cancel CancelFun
	}{
		{
			name:   "grpc",
			server: ServeGRPC,
		},
		{
			name:   "rest",
			server: ServeHttp,
		},
	}

	for _, server := range servers {
		cancelFunChan := make(chan CancelFun, 1)
		go server.server(cancelFunChan)
		server.cancel = <-cancelFunChan
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c

	log.Printf("Stopping service with sig %v", sig)

	for i := len(servers) - 1; i >= 0; i-- {
		server := servers[i]
		if err := server.cancel; err != nil {
			log.Printf("Stop server %v error %v", server.name, err())
		}
	}
}
